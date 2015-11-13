package uncertainty

import (
	"github.com/turing-complete/laboratory/src/internal/config"
)

type Direct struct {
	base
}

func NewDirect(reference []float64, config *config.Uncertainty) (*Direct, error) {
	base, err := newBase(reference, config)
	if err != nil {
		return nil, err
	}
	return &Direct{base}, nil
}

func (self *Direct) Transform(z []float64) []float64 {
	outcome := make([]float64, self.nt)
	copy(outcome, self.lower)
	for i, tid := range self.tasks {
		outcome[tid] += z[i] * (self.upper[tid] - self.lower[tid])
	}
	return outcome
}
