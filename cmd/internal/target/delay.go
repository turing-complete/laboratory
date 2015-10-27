package target

import (
	"github.com/ready-steady/adapt"
	"github.com/simulated-reality/laboratory/cmd/internal/config"
	"github.com/simulated-reality/laboratory/cmd/internal/problem"
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

func (t *delay) String() string {
	return String(t)
}

func (t *delay) Dimensions() (uint, uint) {
	return uint(t.problem.Model.Len()), 2
}

func (t *delay) Compute(node []float64, value []float64) {
	s, m := t.problem.System, t.problem.Model

	value[0] = s.ComputeSchedule(m.Transform(node)).Span
	value[1] = value[0] * value[0]
}

func (t *delay) Monitor(progress *adapt.Progress) {
	if t.config.Verbose {
		Monitor(t, progress)
	}
}

func (t *delay) Score(location *adapt.Location,
	progress *adapt.Progress) float64 {

	return Score(t, t.config, location, progress)
}
