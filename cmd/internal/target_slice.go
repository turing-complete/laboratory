package internal

import (
	"fmt"
	"math"
	"reflect"
	"sync/atomic"
	"unsafe"

	"camlistore.org/pkg/lru"
	"github.com/ready-steady/simulation/power"
	"github.com/ready-steady/simulation/temperature"

	"../../pkg/pool"
)

type sliceTarget struct {
	problem *Problem

	sc uint32
	ec uint32

	power       *power.Power
	temperature *temperature.Temperature

	cache *lru.Cache
	pool  *pool.Pool
}

type sliceData struct {
	P []float64
	S []float64
}

func newSliceTarget(p *Problem) (Target, error) {
	const (
		cacheCapacity = 1000
		poolCapacity  = 100
	)

	c := p.config

	power, err := power.New(p.platform, p.application, c.TempAnalysis.TimeStep)
	if err != nil {
		return nil, err
	}

	temperature, err := temperature.New(temperature.Config(c.TempAnalysis))
	if err != nil {
		return nil, err
	}

	cc, sc := p.cc, uint32(p.schedule.Span/c.TempAnalysis.TimeStep)
	nc := temperature.Nodes

	target := &sliceTarget{
		problem: p,

		sc: sc,

		power:       power,
		temperature: temperature,

		cache: lru.New(cacheCapacity),
		pool: pool.New(poolCapacity, func() interface{} {
			return &sliceData{
				P: make([]float64, cc*sc),
				S: make([]float64, nc*sc),
			}
		}),
	}

	return target, nil
}

func (t *sliceTarget) InputsOutputs() (uint32, uint32) {
	return 1 + t.problem.zc, uint32(len(t.problem.config.CoreIndex)) // +1 for time
}

func (t *sliceTarget) String() string {
	ic, oc := t.InputsOutputs()
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", ic, oc)
}

func (t *sliceTarget) Evaluate(node, value []float64, index []uint64) {
	p := t.problem

	cc, sc := p.cc, t.sc

	var Q []float64
	var key string

	if index != nil {
		key = makeString(index[1:]) // +1 for time
		if result, ok := t.cache.Get(key); ok {
			Q = result.([]float64)
		}
	}

	if Q == nil {
		data := t.pool.Get().(*sliceData)

		u := p.transform(node[1:]) // +1 for time
		t.power.Compute(p.time.Recompute(p.schedule, u), data.P, sc)

		Q = make([]float64, cc*sc)
		t.temperature.ComputeTransient(data.P, Q, data.S, sc)

		t.pool.Put(data)

		if index != nil {
			t.cache.Add(key, Q)
		}

		atomic.AddUint32(&t.ec, 1)
	}

	sid := node[0] * float64(sc-1)
	lid, rid := uint32(math.Floor(sid)), uint32(math.Ceil(sid))

	if lid == rid {
		for i, cid := range p.config.CoreIndex {
			value[i] = Q[lid*cc+uint32(cid)]
		}
	} else {
		fraction := (sid - float64(lid)) / (float64(rid) - float64(lid))
		for i, cid := range p.config.CoreIndex {
			left, right := Q[lid*cc+uint32(cid)], Q[rid*cc+uint32(cid)]
			value[i] = fraction*(right-left) + left
		}
	}
}

func (t *sliceTarget) Progress(level uint8, activeNodes, totalNodes uint32) {
	if level == 0 {
		t.problem.Printf("%10s %15s %15s %15s\n",
			"Level", "Passive Nodes", "Evaluations", "Active Nodes")
	}

	passiveNodes := totalNodes - activeNodes
	t.problem.Printf("%10d %15d %15d %15d\n", level, passiveNodes, t.ec, activeNodes)
}

func makeString(index []uint64) string {
	const (
		sizeOfUInt64 = 8
	)

	sliceHeader := *(*reflect.SliceHeader)(unsafe.Pointer(&index))

	stringHeader := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sizeOfUInt64 * len(index),
	}

	return *(*string)(unsafe.Pointer(&stringHeader))
}
