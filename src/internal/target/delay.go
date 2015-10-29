package target

import (
	"github.com/ready-steady/adapt"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/problem"
)

type delay struct {
	problem *problem.Problem
	config  *config.Target
}

func newDelay(p *problem.Problem, c *config.Target) *delay {
	return &delay{
		problem: p,
		config:  c,
	}
}

func (t *delay) Dimensions() (uint, uint) {
	return uint(t.problem.Uncertainty.Len()), 2
}

func (t *delay) Compute(node []float64, value []float64) {
	s, u := t.problem.System, t.problem.Uncertainty

	value[0] = s.ComputeSchedule(u.Transform(node)).Span
	value[1] = value[0] * value[0]
}

func (t *delay) Monitor(progress *adapt.Progress) {
	if t.config.Verbose {
		monitor(t, progress)
	}
}

func (t *delay) Score(location *adapt.Location, progress *adapt.Progress) float64 {
	return score(t, t.config, location, progress)
}

func (t *delay) String() string {
	return display(t)
}
