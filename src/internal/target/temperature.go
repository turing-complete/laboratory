package target

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type temperature struct {
	base
	time  uncertainty.Uncertainty
	power uncertainty.Uncertainty
}

func newTemperature(system *system.System, config *config.Target) (*temperature, error) {
	base, err := newBase(system, config)
	if err != nil {
		return nil, err
	}

	time, err := uncertainty.New(system.ReferenceTime(), &config.Uncertainty)
	if err != nil {
		return nil, err
	}
	power, err := uncertainty.New(system.ReferencePower(), &config.Uncertainty)
	if err != nil {
		return nil, err
	}
	base.ni = uint(time.Parameters() + power.Parameters())
	base.no = 2 * 1

	return &temperature{base: base, time: time, power: power}, nil
}

func (self *temperature) Compute(node, value []float64) {
	const (
		ε = 1e-10
	)

	nt, np := self.time.Parameters(), self.power.Parameters()

	time := self.time.Transform(node[:nt])
	power := self.time.Transform(node[nt : nt+np])
	schedule := self.system.ComputeSchedule(time)

	P, ΔT := self.system.PartitionPower(power, schedule, ε)
	Q := self.system.ComputeTemperature(P, ΔT)

	max := Q[0]
	for _, q := range Q {
		if q > max {
			max = q
		}
	}

	value[0] = max
	value[1] = max * max
}
