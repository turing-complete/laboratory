package target

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type energy struct {
	base
	uncertainty.Uncertainty
}

func newEnergy(system *system.System, uncertainty *uncertainty.Uncertainty,
	config *config.Target) (*energy, error) {

	ni, _ := uncertainty.Dimensions()
	base, err := newBase(system, config, ni, 2)
	if err != nil {
		return nil, err
	}
	return &energy{base: base, Uncertainty: *uncertainty}, nil
}

func (self *energy) Dimensions() (uint, uint) {
	return self.base.Dimensions()
}

func (self *energy) Compute(node, value []float64) {
	nit, _ := self.Time.Dimensions()
	nip, _ := self.Power.Dimensions()

	time := self.Time.Forward(node[:nit])
	power := self.Power.Forward(node[nit : nit+nip])

	value[0] = 0
	for i := range time {
		value[0] += time[i] * power[i]
	}
	value[1] = value[0] * value[0]
}
