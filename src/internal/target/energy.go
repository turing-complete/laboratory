package target

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type energy struct {
	base
	time  uncertainty.Parameter
	power uncertainty.Parameter
}

func newEnergy(system *system.System, uncertainty *uncertainty.Uncertainty,
	config *config.Target) (*energy, error) {

	base, err := newBase(system, config)
	if err != nil {
		return nil, err
	}

	base.ni, _ = uncertainty.Dimensions()
	base.no = 2 * 1

	return &energy{
		base:  base,
		time:  uncertainty.Time,
		power: uncertainty.Power,
	}, nil
}

func (self *energy) Compute(node, value []float64) {
	nt, _ := self.time.Dimensions()
	np, _ := self.power.Dimensions()

	time := self.time.Forward(node[:nt])
	power := self.power.Forward(node[nt : nt+np])

	value[0] = 0
	for i := range time {
		value[0] += time[i] * power[i]
	}
	value[1] = value[0] * value[0]
}
