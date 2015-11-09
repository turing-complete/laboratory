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
	base, err := newBase(system, config)
	if err != nil {
		return nil, err
	}

	base.uncertainty, err = uncertainty.New(system, system.ReferenceTime(), &config.Uncertainty)
	if err != nil {
		return nil, err
	}
	base.ni, _ = base.uncertainty.Dimensions()
	base.no = 2

	return &energy{base}, nil
}

func (self *energy) Compute(node, value []float64) {
	schedule := self.system.ComputeSchedule(self.uncertainty.Transform(node))
	time, power := self.system.ComputeTime(schedule), self.system.DistributePower(schedule)

	value[0] = 0
	for i := range time {
		value[0] += time[i] * power[i]
	}
	value[1] = value[0] * value[0]
}
