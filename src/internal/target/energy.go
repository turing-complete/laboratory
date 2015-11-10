package target

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type energy struct {
	base
	time  uncertainty.Uncertainty
	power uncertainty.Uncertainty
}

func newEnergy(system *system.System, config *config.Target) (*energy, error) {
	base, err := newBase(system, config)
	if err != nil {
		return nil, err
	}

	time, err := uncertainty.New(system, system.ReferenceTime(), &config.Uncertainty)
	if err != nil {
		return nil, err
	}
	power, err := uncertainty.New(system, system.ReferencePower(), &config.Uncertainty)
	if err != nil {
		return nil, err
	}
	base.ni = uint(time.Len() + power.Len())
	base.no = 2 * 1

	return &energy{base: base, time: time, power: power}, nil
}

func (self *energy) Compute(node, value []float64) {
	nt, np := self.time.Len(), self.power.Len()

	time := self.time.Transform(node[:nt])
	power := self.power.Transform(node[nt : nt+np])

	value[0] = 0
	for i := range time {
		value[0] += time[i] * power[i]
	}
	value[1] = value[0] * value[0]
}
