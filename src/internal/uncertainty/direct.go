package uncertainty

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

type direct struct {
	base
}

func newDirect(system *system.System, reference []float64,
	config *config.Uncertainty) (*direct, error) {

	base, err := newBase(system, reference, config)
	if err != nil {
		return nil, err
	}

	return &direct{base}, nil
}

func (d *direct) Transform(z []float64) []float64 {
	duration := make([]float64, d.nt)
	copy(duration, d.reference)
	for i, tid := range d.taskIndex {
		duration[tid] += z[i] * d.deviation[i]
	}
	return duration
}
