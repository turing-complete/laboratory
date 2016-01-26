package uncertainty

import (
	"fmt"

	"github.com/turing-complete/laboratory/src/internal/config"
)

type epistemic struct {
	base
}

func newEpistemic(reference []float64, config *config.Parameter) (*epistemic, error) {
	base, err := newBase(reference, config)
	if err != nil {
		return nil, err
	}
	return &epistemic{base}, nil
}

func (self *epistemic) Dimensions() (uint, uint) {
	return self.nu, self.nt
}

func (self *epistemic) Forward(ω []float64) []float64 {
	z := make([]float64, self.nu)
	for i, tid := range self.tasks {
		z[i] = (ω[tid] - self.lower[tid]) / (self.upper[tid] - self.lower[tid])
	}
	return z
}

func (self *epistemic) Inverse(z []float64) []float64 {
	ω := make([]float64, self.nt)
	copy(ω, self.lower)
	for i, tid := range self.tasks {
		ω[tid] += z[i] * (self.upper[tid] - self.lower[tid])
	}
	return ω
}

func (self *epistemic) String() string {
	return fmt.Sprintf(`{"dimensions": %d}`, self.nu)
}
