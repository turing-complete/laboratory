package internal

import (
	"fmt"

	"github.com/ready-steady/linear/matrix"
	"github.com/ready-steady/probability/gaussian"

	"../../pkg/solver"
)

type energyTarget struct {
	problem *Problem

	ic uint32 // inputs
}

func newEnergyTarget(p *Problem) (Target, error) {
	return &energyTarget{problem: p, ic: p.zc}, nil
}

func (t *energyTarget) String() string {
	return fmt.Sprintf("Target{inputs: %d, outputs: 1}", t.ic)
}

func (t *energyTarget) InputsOutputs() (uint32, uint32) {
	return t.ic, 1
}

func (t *energyTarget) Serve(jobs <-chan solver.Job) {
	p := t.problem

	cores, tasks := p.platform.Cores, p.application.Tasks

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

		schedule := p.time.Recompute(p.schedule, d)

		job.Value[0] = 0
		for i := uint32(0); i < tc; i++ {
			job.Value[0] += (schedule.Finish[i] - schedule.Start[i]) *
				cores[uint32(schedule.Mapping[i])].Power[tasks[i].Type]
		}

		job.Done <- solver.Result{}
	}
}
