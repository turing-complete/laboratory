package target

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type temperature struct {
	base
	time  *uncertainty.Parameter
	power *uncertainty.Parameter
}

func newTemperature(system *system.System, uncertainty *uncertainty.Uncertainty,
	config *config.Target) (*temperature, error) {

	base, err := newBase(system, config)
	if err != nil {
		return nil, err
	}

	base.ni = uncertainty.Time.Dimensions() + uncertainty.Power.Dimensions()
	base.no = 2 * 1

	return &temperature{
		base:  base,
		time:  uncertainty.Time,
		power: uncertainty.Power,
	}, nil
}

func (self *temperature) Compute(node, value []float64) {
	const (
		ε = 1e-10
	)

	nt, np := self.time.Dimensions(), self.power.Dimensions()

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
