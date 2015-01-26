package internal

import (
	"fmt"

	"github.com/ready-steady/linear/matrix"
	"github.com/ready-steady/probability/gaussian"

	"../../pkg/solver"
)

type spanTarget struct {
	problem *Problem

	ic uint32 // inputs
}

func newSpanTarget(p *Problem) (Target, error) {
	return &spanTarget{problem: p, ic: p.zc}, nil
}

func (t *spanTarget) String() string {
	return fmt.Sprintf("Target{inputs: %d, outputs: 1}", t.ic)
}

func (t *spanTarget) InputsOutputs() (uint32, uint32) {
	return t.ic, 1
}

func (t *spanTarget) Serve(jobs <-chan solver.Job) {
	p := t.problem
	c := &p.config

	zc, uc, tc := p.zc, p.uc, p.tc

	g := gaussian.New(0, 1)
	m := p.marginals

	z := make([]float64, zc)
	u := make([]float64, uc)
	d := make([]float64, tc)

	for job := range jobs {
		// Independent uniform to independent Gaussian
		for i := uint32(0); i < zc; i++ {
			z[i] = g.InvCDF(processNode(job.Node[i]))
		}

		// Independent Gaussian to dependent Gaussian
		matrix.Multiply(p.transform, z, u, uc, zc, 1)

		// Dependent Gaussian to dependent uniform to dependent target
		for i, tid := range c.TaskIndex {
			d[tid] = m[i].InvCDF(g.CDF(u[i]))
		}

		job.Value[0] = p.time.Recompute(p.schedule, d).Span

		job.Done <- solver.Result{}
	}
}
