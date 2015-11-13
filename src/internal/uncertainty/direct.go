package uncertainty

import (
	"github.com/turing-complete/laboratory/src/internal/config"
)

type direct struct {
	base
}

func newDirect(reference []float64, config *config.Uncertainty) (*direct, error) {
	base, err := newBase(reference, config)
	if err != nil {
		return nil, err
	}
	return &direct{base}, nil
}

func (self *direct) Transform(z []float64) []float64 {
	outcome := make([]float64, self.nt)
	copy(outcome, self.lower)
	for i, tid := range self.tasks {
		outcome[tid] += z[i] * (self.upper[tid] - self.lower[tid])
	}
	return outcome
}
