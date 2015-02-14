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

	pc uint
	sc uint
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

	c := &p.Config

	power, err := power.New(p.platform, p.application, c.TempAnalysis.TimeStep)
	if err != nil {
		return nil, err
	}

	temperature, err := temperature.New(temperature.Config(c.TempAnalysis))
	if err != nil {
		return nil, err
	}

	cc, sc := p.cc, uint(p.schedule.Span/c.TempAnalysis.TimeStep)
	nc := temperature.Nodes

	target := &sliceTarget{
		problem: p,

		pc: 1, // +1 for time
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

func (t *sliceTarget) Inputs() uint {
	return t.pc + t.problem.zc
}

func (t *sliceTarget) Outputs() uint {
	return uint(len(t.problem.Config.CoreIndex))
}

func (t *sliceTarget) Pseudos() uint {
	return t.pc
}

func (t *sliceTarget) String() string {
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", t.Inputs(), t.Outputs())
}

func (t *sliceTarget) Evaluate(node, value []float64, index []uint64) {
	p := t.problem

	cc, pc, sc := p.cc, t.pc, t.sc

	var Q []float64
	var key string

	if index != nil {
		key = makeString(index[pc:])
		if result, ok := t.cache.Get(key); ok {
			Q = result.([]float64)
		}
	}

	if Q == nil {
		data := t.pool.Get().(*sliceData)

		u := p.transform(node[pc:])
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
	lid, rid := uint(math.Floor(sid)), uint(math.Ceil(sid))

	if lid == rid {
		for i, cid := range p.Config.CoreIndex {
			value[i] = Q[lid*cc+cid]
		}
	} else {
		fraction := (sid - float64(lid)) / (float64(rid) - float64(lid))
		for i, cid := range p.Config.CoreIndex {
			left, right := Q[lid*cc+cid], Q[rid*cc+cid]
			value[i] = fraction*(right-left) + left
		}
	}
}

func (t *sliceTarget) Progress(level uint32, activeNodes, totalNodes uint) {
	if level == 0 {
		fmt.Printf("%10s %15s %15s %15s\n", "Level", "Passive Nodes", "Evaluations", "Active Nodes")
	}

	passiveNodes := totalNodes - activeNodes
	fmt.Printf("%10d %15d %15d %15d\n", level, passiveNodes, t.ec, activeNodes)
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
