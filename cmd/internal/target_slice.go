package internal

import (
	"fmt"
	"reflect"
	"sync/atomic"
	"unsafe"

	"camlistore.org/pkg/lru"
	"github.com/ready-steady/numeric/interpolation/spline"
	"github.com/ready-steady/simulation/power"
	"github.com/ready-steady/simulation/temperature/numeric"
)

type sliceTarget struct {
	problem *Problem

	np uint
	ne uint32

	interval    []float64
	power       *power.Power
	temperature *numeric.Temperature

	cache *lru.Cache
}

func newSliceTarget(p *Problem) (Target, error) {
	const (
		cacheCapacity = 1000
	)

	c := &p.Config

	power := power.New(p.platform, p.application)
	temperature, err := numeric.New((*numeric.Config)(&c.Temperature))
	if err != nil {
		return nil, err
	}

	target := &sliceTarget{
		problem: p,

		np: 1, // +1 for time

		interval:    []float64{0, p.schedule.Span},
		power:       power,
		temperature: temperature,

		cache: lru.New(cacheCapacity),
	}

	return target, nil
}

func (t *sliceTarget) Inputs() uint {
	return t.np + t.problem.nz
}

func (t *sliceTarget) Outputs() uint {
	return uint(len(t.problem.Config.CoreIndex))
}

func (t *sliceTarget) Pseudos() uint {
	return t.np
}

func (t *sliceTarget) String() string {
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", t.Inputs(), t.Outputs())
}

func (t *sliceTarget) Evaluate(node, value []float64, index []uint64) {
	p := t.problem

	np := t.np

	var interpolant *spline.Cubic
	var key string

	if index != nil {
		key = makeString(index[np:])
		if result, ok := t.cache.Get(key); ok {
			interpolant = result.(*spline.Cubic)
		}
	}

	if interpolant == nil {
		schedule := p.time.Recompute(p.schedule, p.transform(node[np:]))
		Q, time, err := t.temperature.Compute(t.power.Process(schedule), t.interval)
		if err != nil {
			panic("cannot compute a temperature profile")
		}

		interpolant = spline.NewCubic(time, Q)

		atomic.AddUint32(&t.ne, 1)

		if index != nil {
			t.cache.Add(key, interpolant)
		}
	}

	Q := interpolant.Compute([]float64{node[0] * t.interval[1]})
	for i, cid := range p.Config.CoreIndex {
		value[i] = Q[cid]
	}
}

func (t *sliceTarget) Progress(level uint32, na, nt uint) {
	if level == 0 {
		fmt.Printf("%10s %15s %15s %15s\n", "Level", "Passive Nodes", "Evaluations", "Active Nodes")
	}
	fmt.Printf("%10d %15d %15d %15d\n", level, nt-na, t.ne, na)
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
