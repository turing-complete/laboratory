package uncertainty

import (
	"errors"
	"fmt"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/support"
	"github.com/turing-complete/laboratory/src/internal/system"
)

type base struct {
	taskIndex []uint
	reference []float64
	deviation []float64

	nt uint
	nu uint
	nz uint
}

func newBase(system *system.System, reference []float64,
	config *config.Uncertainty) (*base, error) {

	nt := uint(system.Application.Len())
	if nt != uint(len(reference)) {
		return nil, errors.New("the length of the reference is invalid")
	}

	taskIndex, err := support.ParseNaturalIndex(config.TaskIndex, 0, nt-1)
	if err != nil {
		return nil, err
	}

	nu := uint(len(taskIndex))

	deviation := make([]float64, nu)
	for i, tid := range taskIndex {
		deviation[i] = config.MaxDeviation * reference[tid]
	}

	return &base{
		taskIndex: taskIndex,
		reference: reference,
		deviation: deviation,

		nt: nt,
		nu: nu,
		nz: nu,
	}, nil
}

func (b *base) Len() int {
	return int(b.nz)
}

func (b *base) String() string {
	return fmt.Sprintf(`{"parameters": %d, "variables": %d}`, b.nu, b.nz)
}
