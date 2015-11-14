package uncertainty

import (
	"fmt"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/support"
)

type direct struct {
	tasks []uint
	lower []float64
	upper []float64

	nt uint
	nu uint
}

func newDirect(reference []float64, config *config.Uncertainty) (*direct, error) {
	nt := uint(len(reference))

	tasks, err := support.ParseNaturalIndex(config.Tasks, 0, nt-1)
	if err != nil {
		return nil, err
	}

	nu := uint(len(tasks))

	lower := make([]float64, nt)
	upper := make([]float64, nt)

	copy(lower, reference)
	copy(upper, reference)

	for _, tid := range tasks {
		upper[tid] *= (1.0 + config.Deviation)
	}

	return &direct{
		tasks: tasks,
		lower: lower,
		upper: upper,

		nt: nt,
		nu: nu,
	}, nil
}

func (self *direct) Dimensions() uint {
	return self.nu
}

func (self *direct) String() string {
	return fmt.Sprintf(`{"dimensions": %d}`, self.nu)
}

func (self *direct) Transform(z []float64) []float64 {
	outcome := make([]float64, self.nt)
	copy(outcome, self.lower)
	for i, tid := range self.tasks {
		outcome[tid] += z[i] * (self.upper[tid] - self.lower[tid])
	}
	return outcome
}
