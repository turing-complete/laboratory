package quantity

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/support"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type energy struct {
	base
}

func newEnergy(system *system.System, uncertainty uncertainty.Uncertainty,
	config *config.Quantity) (*energy, error) {

	ni, _ := uncertainty.Mapping()
	base, err := newBase(system, uncertainty, config, ni, 1)
	if err != nil {
		return nil, err
	}
	return &energy{base}, nil
}

func (self *energy) Compute(node, value []float64) {
	P := self.system.ComputeDynamicPower(self.system.ComputeSchedule(self.Backward(node)))
	self.system.ComputeTemperatureUpdatePower(P)
	value[0] = support.Sum(P) * self.system.TimeStep()
}
