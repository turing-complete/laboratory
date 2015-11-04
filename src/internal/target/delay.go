package target

import (
	"github.com/ready-steady/adapt"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type delay struct {
	system *system.System
	config *config.Target

	uncertainty uncertainty.Uncertainty
}

func newDelay(system *system.System, config *config.Target) (*delay, error) {
	uncertainty, err := uncertainty.New(system, &config.Uncertainty)
	if err != nil {
		return nil, err
	}

	return &delay{
		system: system,
		config: config,

		uncertainty: uncertainty,
	}, nil
}

func (t *delay) Dimensions() (uint, uint) {
	return uint(t.uncertainty.Len()), 2
}

func (t *delay) Compute(node []float64, value []float64) {
	value[0] = t.system.ComputeSchedule(t.uncertainty.Transform(node)).Span
	value[1] = value[0] * value[0]
}

func (t *delay) Monitor(progress *adapt.Progress) {
	monitor(t, progress)
}

func (t *delay) Score(location *adapt.Location, progress *adapt.Progress) float64 {
	return score(t, t.config, location, progress)
}

func (t *delay) String() string {
	return display(t)
}
