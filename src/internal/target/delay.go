package target

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type delay struct {
	base
	time uncertainty.Uncertainty
}

func newDelay(system *system.System, config *config.Target) (*delay, error) {
	base, err := newBase(system, config)
	if err != nil {
		return nil, err
	}

	time, err := uncertainty.New(system.ReferenceTime(), &config.Uncertainty)
	if err != nil {
		return nil, err
	}
	base.ni = uint(time.Parameters())
	base.no = 2 * 1

	return &delay{base: base, time: time}, nil
}

func (self *delay) Compute(node []float64, value []float64) {
	value[0] = self.system.ComputeSchedule(self.time.Transform(node)).Span
	value[1] = value[0] * value[0]
}
