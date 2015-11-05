package target

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type energy struct {
	base
}

func newEnergy(system *system.System, config *config.Target) (*energy, error) {
	uncertainty, err := uncertainty.New(system, &config.Uncertainty)
	if err != nil {
		return nil, err
	}

	return &energy{
		base: base{
			system:      system,
			config:      config,
			uncertainty: uncertainty,

			ni: uint(uncertainty.Len()),
			no: 2,
		},
	}, nil
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
