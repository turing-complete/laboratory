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
	uncertainty, err := uncertainty.New(system, &config.Uncertainty)
	if err != nil {
		return nil, err
	}

	return &delay{
		base: base{
			system:      system,
			config:      config,
			uncertainty: uncertainty,

			ni: uint(uncertainty.Len()),
			no: 2,
		},
	}, nil
}

func (t *delay) Compute(node []float64, value []float64) {
	value[0] = t.system.ComputeSchedule(t.uncertainty.Transform(node)).Span
	value[1] = value[0] * value[0]
}
