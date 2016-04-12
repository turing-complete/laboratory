package target

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type delay struct {
	base
	uncertainty.Parameter
}

func newDelay(system *system.System, uncertainty *uncertainty.Uncertainty,
	config *config.Target) (*delay, error) {

	ni, _ := uncertainty.Time.Mapping()
	base, err := newBase(system, config, ni, 1)
	if err != nil {
		return nil, err
	}
	return &delay{base: base, Parameter: uncertainty.Time}, nil
}

func (self *delay) Compute(node []float64, value []float64) {
	value[0] = self.system.ComputeSchedule(self.Inverse(node)).Span
}
