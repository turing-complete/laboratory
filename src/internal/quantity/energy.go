package quantity

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type energy struct {
	base

	power []float64
}

func newEnergy(system *system.System, uncertainty uncertainty.Uncertainty,
	config *config.Quantity) (*energy, error) {

	ni, _ := uncertainty.Mapping()
	base, err := newBase(system, uncertainty, config, ni, 1)
	if err != nil {
		return nil, err
	}
	return &energy{
		base:  base,
		power: system.ReferencePower(),
	}, nil
}

func (self *energy) Compute(node, value []float64) {
	time := self.Backward(node)

	value[0] = 0.0
	for i, power := range self.power {
		value[0] += time[i] * power
	}
}
