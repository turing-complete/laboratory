package target

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type delay struct {
	base
}

func newDelay(system *system.System, config *config.Target) (*delay, error) {
	base, err := newBase(system, config)
	if err != nil {
		return nil, err
	}

	base.uncertainty, err = uncertainty.New(system, system.ReferenceTime(), &config.Uncertainty)
	if err != nil {
		return nil, err
	}
	base.ni, _ = base.uncertainty.Dimensions()
	base.no = 2

	return &delay{base}, nil
}

func (self *delay) Compute(node []float64, value []float64) {
	value[0] = self.system.ComputeSchedule(self.uncertainty.Transform(node)).Span
	value[1] = value[0] * value[0]
}
