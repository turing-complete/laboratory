package quantity

import (
	"math"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type temperature struct {
	base

	power []float64
}

func newTemperature(system *system.System, uncertainty uncertainty.Uncertainty,
	config *config.Quantity) (*temperature, error) {

	ni, _ := uncertainty.Mapping()
	base, err := newBase(system, uncertainty, config, ni, 1)
	if err != nil {
		return nil, err
	}
	return &temperature{
		base:  base,
		power: system.ReferencePower(),
	}, nil
}

func (self *temperature) Compute(node, value []float64) {
	const (
		ε = 1e-10
	)

	schedule := self.system.ComputeSchedule(self.Backward(node))
	P, ΔT := self.system.PartitionPower(self.power, schedule, ε)
	Q := self.system.ComputeTemperature(P, ΔT)

	value[0] = Q[0]
	for _, q := range Q {
		value[0] = math.Max(value[0], q)
	}
}
