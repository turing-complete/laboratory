package main

import (
	"fmt"

	"github.com/ready-steady/linal/matrix"
	"github.com/ready-steady/prob/gaussian"
)

type spanTarget struct {
	problem *problem

	ic uint32 // inputs
}

func newSpanTarget(p *problem) (target, error) {
	return &spanTarget{problem: p, ic: p.zc}, nil
}

func (t *spanTarget) String() string {
	return fmt.Sprintf("Target{inputs: %d, outputs: 1}", t.ic)
}

func (t *spanTarget) InputsOutputs() (uint32, uint32) {
	return t.ic, 1
}

func (t *spanTarget) Serve(jobs <-chan job) {
	p := t.problem
	c := &p.config

	uc, zc := p.uc, p.zc

	g := gaussian.New(0, 1)
	m := p.marginals

	z := make([]float64, zc)
	u := make([]float64, uc)
	d := make([]float64, p.tc)

	for job := range jobs {
		// Independent uniform to independent Gaussian
		for i := uint32(0); i < zc; i++ {
			z[i] = g.InvCDF(job.node[i])
		}

		// Independent Gaussian to dependent Gaussian
		matrix.Multiply(p.transform, z, u, uc, zc, 1)

		// Dependent Gaussian to dependent uniform to dependent target
		for i, tid := range c.TaskIndex {
			d[tid] = m[i].InvCDF(g.CDF(u[i]))
		}

		job.value[0] = p.time.Recompute(p.schedule, d).Span

		job.done <- result{}
	}
}
