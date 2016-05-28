package quantity

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type energy struct {
	base

	Δt float64
}

func newEnergy(system *system.System, uncertainty uncertainty.Uncertainty,
	config *config.Quantity) (*energy, error) {

	ni, _ := uncertainty.Mapping()
	base, err := newBase(system, uncertainty, config, ni, 1)
	if err != nil {
		return nil, err
	}
	return &energy{
		base: base,
		Δt:   system.TimeStep(),
	}, nil
}

func (self *energy) Compute(node, value []float64) {
	P := self.system.ComputePower(self.system.ComputeSchedule(self.Backward(node)))
	value[0] = 0.0
	for _, p := range P {
		value[0] += p
	}
	value[0] *= self.Δt
}
