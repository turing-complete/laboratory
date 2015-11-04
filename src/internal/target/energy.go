package target

import (
	"github.com/ready-steady/adapt"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type energy struct {
	system *system.System
	config *config.Target

	uncertainty uncertainty.Uncertainty
}

func newEnergy(system *system.System, config *config.Target) (*energy, error) {
	uncertainty, err := uncertainty.New(system, &config.Uncertainty)
	if err != nil {
		return nil, err
	}

	return &energy{
		system: system,
		config: config,

		uncertainty: uncertainty,
	}, nil
}

func (t *energy) Dimensions() (uint, uint) {
	return uint(t.uncertainty.Len()), 2
}

func (t *energy) Compute(node, value []float64) {
	schedule := t.system.ComputeSchedule(t.uncertainty.Transform(node))
	time, power := t.system.ComputeTime(schedule), t.system.DistributePower(schedule)

	value[0] = 0
	for i := range time {
		value[0] += time[i] * power[i]
	}
	value[1] = value[0] * value[0]
}

func (t *energy) Monitor(progress *adapt.Progress) {
	if t.config.Verbose {
		monitor(t, progress)
	}
}

func (t *energy) Score(location *adapt.Location, progress *adapt.Progress) float64 {
	return score(t, t.config, location, progress)
}

func (t *energy) String() string {
	return display(t)
}
