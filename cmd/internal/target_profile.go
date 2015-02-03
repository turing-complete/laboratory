package internal

// #include <string.h>
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/ready-steady/simulation/power"
	"github.com/ready-steady/simulation/temperature"

	"../../pkg/pool"
)

type profileTarget struct {
	problem *Problem

	sc uint32

	power       *power.Power
	temperature *temperature.Temperature

	pool *pool.Pool
}

type profileData struct {
	P []float64
	S []float64
}

func newProfileTarget(p *Problem) (Target, error) {
	const (
		poolCapacity = 100
		MaxUInt16    = ^uint16(0)
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

	if cc*sc > uint32(MaxUInt16) {
		panic("The number of outputs is too large.")
	}

	target := &profileTarget{
		problem: p,

		sc: sc,

		power:       power,
		temperature: temperature,

		pool: pool.New(poolCapacity, func() interface{} {
			return &profileData{
				P: make([]float64, cc*sc),
				S: make([]float64, nc*sc),
			}
		}),
	}

	return target, nil
}

func (t *profileTarget) InputsOutputs() (uint32, uint32) {
	return t.problem.zc, t.sc * t.problem.cc
}

func (t *profileTarget) String() string {
	ic, oc := t.InputsOutputs()
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", ic, oc)
}

func (t *profileTarget) Evaluate(node, value []float64, _ []uint64) {
	p := t.problem

	cc, sc := p.cc, t.sc

	u := p.transform(node)

	data := t.pool.Get().(*profileData)

	// FIXME: Bad, bad, bad!
	C.memset(unsafe.Pointer(&data.P[0]), 0, C.size_t(8*cc*sc))

	t.power.Compute(p.time.Recompute(p.schedule, u), data.P, sc)
	t.temperature.ComputeTransient(data.P, value, data.S, sc)

	t.pool.Put(data)
}

func (t *profileTarget) Progress(level uint8, activeNodes, totalNodes uint32) {
	passiveNodes := totalNodes - activeNodes
	t.problem.Printf("%5d %10d %10d\n", level, passiveNodes, activeNodes)
}
