package uncertainty

import (
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

func newBase(reference []float64, config *config.Parameter) (base, error) {
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
		lower[tid] -= config.Deviation * reference[tid]
		upper[tid] += config.Deviation * reference[tid]
	}

	return base{
		tasks: tasks,
		lower: lower,
		upper: upper,

		nt: nt,
		nu: nu,
	}, nil
}
