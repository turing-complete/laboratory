package target

import (
	"math"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type temperature struct {
	base
}

func newTemperature(system *system.System, uncertainty *uncertainty.Uncertainty,
	config *config.Target) (*temperature, error) {

	ni, _ := uncertainty.Mapping()
	base, err := newBase(system, uncertainty, config, ni, 1)
	if err != nil {
		return nil, err
	}
	return &temperature{base}, nil
}

func (self *temperature) Compute(node, value []float64) {
	const (
		ε = 1e-10
	)

	nt := uint(self.system.Application.Len())
	timePower := self.Inverse(node)

	schedule := self.system.ComputeSchedule(timePower[:nt])
	P, ΔT := self.system.PartitionPower(timePower[nt:], schedule, ε)
	Q := self.system.ComputeTemperature(P, ΔT)

	value[0] = Q[0]
	for _, q := range Q {
		value[0] = math.Max(value[0], q)
	}
}
