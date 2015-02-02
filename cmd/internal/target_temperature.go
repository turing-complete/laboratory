package internal

// #include <string.h>
import "C"

import (
	"fmt"
	"math"
	"sync/atomic"
	"unsafe"

	"github.com/ready-steady/simulation/power"
	"github.com/ready-steady/simulation/temperature"

	"../../pkg/cache"
	"../../pkg/pool"
)

type temperatureTarget struct {
	problem *Problem

	sc uint32
	ec uint32

	power       *power.Power
	temperature *temperature.Temperature

	cache *cache.Cache
	pool  *pool.Pool
}

type temperatureData struct {
	P []float64
	S []float64
}

func newTemperatureTarget(p *Problem) (Target, error) {
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

	target := &temperatureTarget{
		problem: p,

		sc: sc,

		power:       power,
		temperature: temperature,

		cache: cache.New(p.zc, cacheCapacity),
		pool: pool.New(poolCapacity, func() interface{} {
			return &temperatureData{
				P: make([]float64, cc*sc),
				S: make([]float64, nc*sc),
			}
		}),
	}

	return target, nil
}

func (t *temperatureTarget) Evaluate(node, value []float64, index []uint64) {
	p := t.problem

	cc, sc := p.cc, t.sc

	var Q []float64
	var key string

	if index != nil {
		key = t.cache.Key(index[1:]) // +1 for time
		Q = t.cache.Get(key)
	}

	if Q == nil {
		data := t.pool.Get().(*temperatureData)

		// FIXME: Bad, bad, bad!
		C.memset(unsafe.Pointer(&data.P[0]), 0, C.size_t(8*cc*sc))

		u := p.transform(node[1:]) // +1 for time
		t.power.Compute(p.time.Recompute(p.schedule, u), data.P, sc)

		Q = make([]float64, cc*sc)
		t.temperature.ComputeTransient(data.P, Q, data.S, sc)

		t.pool.Put(data)

		if index != nil {
			t.cache.Set(key, Q)
		}

		atomic.AddUint32(&t.ec, 1)
	}

	sid := node[0] * float64(sc-1)
	lid, rid := uint32(math.Floor(sid)), uint32(math.Ceil(sid))

	if lid == rid {
		for i := uint32(0); i < cc; i++ {
			value[i] = Q[lid*cc+i]
		}
	} else {
		fraction := (sid - float64(lid)) / (float64(rid) - float64(lid))
		for i := uint32(0); i < cc; i++ {
			left, right := Q[lid*cc+i], Q[rid*cc+i]
			value[i] = fraction*(right-left) + left
		}
	}
}

func (t *temperatureTarget) InputsOutputs() (uint32, uint32) {
	return 1 + t.problem.zc, t.problem.cc // +1 for time
}

func (t *temperatureTarget) Evaluations() uint32 {
	return t.ec
}

func (t *temperatureTarget) String() string {
	ic, oc := t.InputsOutputs()
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", ic, oc)
}
