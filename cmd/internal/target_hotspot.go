package internal

// #include <string.h>
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/ready-steady/linear/matrix"
	"github.com/ready-steady/probability/gaussian"
	"github.com/ready-steady/simulation/power"
	"github.com/ready-steady/simulation/temperature"

	"../../pkg/solver"
)

type hotspotTarget struct {
	problem *Problem

	ic uint32 // inputs
	oc uint32 // outputs
	sc uint32 // steps

	power       *power.Power
	temperature *temperature.Temperature
}

func newHotspotTarget(p *Problem) (Target, error) {
	c := &p.config

	power, err := power.New(p.platform, p.application, c.TempAnalysis.TimeStep)
	if err != nil {
		return nil, err
	}

	temperature, err := temperature.New(temperature.Config(c.TempAnalysis))
	if err != nil {
		return nil, err
	}

	target := &hotspotTarget{
		problem: p,

		ic: 1 + p.zc, // +1 for time
		oc: uint32(len(c.CoreIndex)),
		sc: uint32(p.schedule.Span / c.TempAnalysis.TimeStep),

		power:       power,
		temperature: temperature,
	}

	return target, nil
}

func (t *hotspotTarget) String() string {
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", t.ic, t.oc)
}

func (t *hotspotTarget) InputsOutputs() (uint32, uint32) {
	return t.ic, t.oc
}

func (t *hotspotTarget) Serve(jobs <-chan solver.Job) {
	p := t.problem
	c := &p.config

	zc, uc, cc, tc, sc := p.zc, p.uc, p.cc, p.tc, t.sc

	g := gaussian.New(0, 1)
	m := p.marginals

	P := make([]float64, cc*sc)
	S := make([]float64, t.temperature.Nodes*sc)
	Q := make([]float64, cc*sc)

	z := make([]float64, zc)
	u := make([]float64, uc)
	d := make([]float64, tc)

	for job := range jobs {
		// Independent uniform to independent Gaussian
		for i := uint32(0); i < zc; i++ {
			z[i] = g.InvCDF(processNode(job.Node[1+i])) // +1 for time
		}

		// Independent Gaussian to dependent Gaussian
		matrix.Multiply(p.transform, z, u, uc, zc, 1)

		// Dependent Gaussian to dependent uniform to dependent target
		for i, tid := range c.TaskIndex {
			d[tid] = m[i].InvCDF(g.CDF(u[i]))
		}

		// FIXME: Bad, bad, bad!
		C.memset(unsafe.Pointer(&P[0]), 0, C.size_t(8*cc*sc))

		t.power.Compute(p.time.Recompute(p.schedule, d), P, sc)
		t.temperature.ComputeTransient(P, Q, S, sc)

		value := 0.0
		for i := range Q {
			if Q[i] > value {
				value = Q[i]
			}
		}
		job.Value[0] = value

		job.Done <- solver.Result{}
	}
}
