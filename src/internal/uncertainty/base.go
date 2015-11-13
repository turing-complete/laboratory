package uncertainty

import (
	"fmt"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/support"
)

type base struct {
	tasks []uint
	lower []float64
	upper []float64

	nt uint
	nu uint
}

func newBase(reference []float64, config *config.Uncertainty) (base, error) {
	nt := uint(len(reference))

	tasks, err := support.ParseNaturalIndex(config.Tasks, 0, nt-1)
	if err != nil {
		return base{}, err
	}

	nu := uint(len(tasks))

	lower := make([]float64, nt)
	upper := make([]float64, nt)

	copy(lower, reference)
	copy(upper, reference)

	for _, tid := range tasks {
		upper[tid] *= (1.0 + config.Deviation)
	}

	return base{
		tasks: tasks,
		lower: lower,
		upper: upper,

		nt: nt,
		nu: nu,
	}, nil
}

func (self *base) Parameters() uint {
	return self.nu
}

func (self *base) String() string {
	return fmt.Sprintf(`{"parameters": %d}`, self.nu)
}
