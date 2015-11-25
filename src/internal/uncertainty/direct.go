package uncertainty

import (
	"fmt"

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

func (self *direct) Dimensions() (uint, uint) {
	return self.nu, self.nt
}

func (self *direct) Forward(z []float64) []float64 {
	ω := make([]float64, self.nt)
	copy(ω, self.lower)
	for i, tid := range self.tasks {
		ω[tid] += z[i] * (self.upper[tid] - self.lower[tid])
	}
	return ω
}

func (self *direct) Inverse(ω []float64) []float64 {
	z := make([]float64, self.nu)
	for i, tid := range self.tasks {
		z[i] = (ω[tid] - self.lower[tid]) / (self.upper[tid] - self.lower[tid])
	}
	return z
}

func (self *direct) String() string {
	return fmt.Sprintf(`{"dimensions": %d}`, self.nu)
}
