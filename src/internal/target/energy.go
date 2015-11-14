package target

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type energy struct {
	base
	time  *uncertainty.Parameter
	power *uncertainty.Parameter
}

func newEnergy(system *system.System, uncertainty *uncertainty.Uncertainty,
	config *config.Target) (*energy, error) {

	base, err := newBase(system, config)
	if err != nil {
		return nil, err
	}

	base.ni = uint(uncertainty.Time.Len() + uncertainty.Power.Len())
	base.no = 2 * 1

	return &energy{
		base:  base,
		time:  uncertainty.Time,
		power: uncertainty.Power,
	}, nil
}

func (self *energy) Compute(node, value []float64) {
	nt, np := uint(self.time.Len()), uint(self.power.Len())

	time := self.time.Transform(node[:nt])
	power := self.power.Transform(node[nt : nt+np])

	value[0] = 0
	for i := range time {
		value[0] += time[i] * power[i]
	}
	value[1] = value[0] * value[0]
}
