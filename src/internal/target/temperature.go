package target

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type temperature struct {
	base
	uncertainty.Uncertainty
}

func newTemperature(system *system.System, uncertainty *uncertainty.Uncertainty,
	config *config.Target) (*temperature, error) {

	ni, _ := uncertainty.Dimensions()
	base, err := newBase(system, config, ni, 1)
	if err != nil {
		return nil, err
	}
	return &temperature{base: base, Uncertainty: *uncertainty}, nil
}

func (self *temperature) Dimensions() (uint, uint) {
	return self.base.Dimensions()
}

func (self *temperature) Compute(node, value []float64) {
	const (
		ε = 1e-10
	)

	nit, _ := self.Time.Dimensions()
	nip, _ := self.Power.Dimensions()

	time := self.Time.Inverse(node[:nit])
	power := self.Power.Inverse(node[nit : nit+nip])
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
}
