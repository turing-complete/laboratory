package target

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type energy struct {
	base
}

func newEnergy(system *system.System, uncertainty *uncertainty.Uncertainty,
	config *config.Target) (*energy, error) {

	ni, _ := uncertainty.Mapping()
	base, err := newBase(system, uncertainty, config, ni, 1)
	if err != nil {
		return nil, err
	}
	return &energy{base}, nil
}

func (self *energy) Compute(node, value []float64) {
	nt := uint(self.system.Application.Len())
	timePower := self.Inverse(node)

	value[0] = 0
	for i := uint(0); i < nt; i++ {
		value[0] += timePower[i] * timePower[nt+i]
	}
}
