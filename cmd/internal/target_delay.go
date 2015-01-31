package internal

import (
	"fmt"

	"github.com/ready-steady/linear/matrix"
	"github.com/ready-steady/probability/gaussian"

	"../../pkg/solver"
)

type delayTarget struct {
	problem *Problem

	ic uint32 // inputs
}

func newDelayTarget(p *Problem) (Target, error) {
	return &delayTarget{problem: p, ic: p.zc}, nil
}

func (t *delayTarget) String() string {
	return fmt.Sprintf("Target{inputs: %d, outputs: 1}", t.ic)
}

func (t *delayTarget) InputsOutputs() (uint32, uint32) {
	return t.ic, 1
}

func (t *delayTarget) Serve(jobs <-chan solver.Job) {
	p := t.problem

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
		for i := uint32(0); i < tc; i++ {
			d[i] = m[i].InvCDF(g.CDF(u[i]))
		}

		job.Value[0] = p.time.Recompute(p.schedule, d).Span

		job.Done <- solver.Result{}
	}
}
