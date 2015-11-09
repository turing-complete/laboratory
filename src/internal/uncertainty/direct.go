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

func (self *direct) Transform(z []float64) []float64 {
	duration := make([]float64, self.nt)
	copy(duration, self.reference)
	for i, tid := range self.taskIndex {
		duration[tid] += z[i] * self.deviation[i]
	}
	return duration
}
