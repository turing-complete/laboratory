package quantity

import (
	"math"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type temperature struct {
	base
}

func newTemperature(system *system.System, uncertainty uncertainty.Uncertainty,
	config *config.Quantity) (*temperature, error) {

	ni, _ := uncertainty.Mapping()
	base, err := newBase(system, uncertainty, config, ni, 1)
	if err != nil {
		return nil, err
	}
	return &temperature{base}, nil
}

func (self *temperature) Compute(node, value []float64) {
	P := self.system.ComputeDynamicPower(self.system.ComputeSchedule(self.Backward(node)))
	Q := self.system.ComputeTemperatureUpdatePower(P)
	value[0] = 0.0
	for _, q := range Q {
		value[0] = math.Max(value[0], q)
	}
}
