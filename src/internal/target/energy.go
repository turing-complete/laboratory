package target

import (
	"github.com/ready-steady/adapt"
	"github.com/simulated-reality/laboratory/src/internal/config"
	"github.com/simulated-reality/laboratory/src/internal/problem"
)

type energy struct {
	problem *problem.Problem
	config  *config.Target
}

func newEnergy(p *problem.Problem, c *config.Target) *energy {
	return &energy{
		problem: p,
		config:  c,
	}
}

func (t *energy) String() string {
	return String(t)
}

func (t *energy) Dimensions() (uint, uint) {
	return uint(t.problem.Uncertainty.Len()), 2
}

func (t *energy) Compute(node, value []float64) {
	s, u := t.problem.System, t.problem.Uncertainty

	schedule := s.ComputeSchedule(u.Transform(node))
	time, power := s.ComputeTime(schedule), s.DistributePower(schedule)

	value[0] = 0
	for i := range time {
		value[0] += time[i] * power[i]
	}

	value[1] = value[0] * value[0]
}

func (t *energy) Monitor(progress *adapt.Progress) {
	if t.config.Verbose {
		Monitor(t, progress)
	}
}

func (t *energy) Score(location *adapt.Location, progress *adapt.Progress) float64 {
	return Score(t, t.config, location, progress)
}
