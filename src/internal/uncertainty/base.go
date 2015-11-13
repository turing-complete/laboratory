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
	config *config.Uncertainty) (base, error) {

	nt := uint(system.Application.Len())
	if nt != uint(len(reference)) {
		return base{}, errors.New("the length of the reference is invalid")
	}

	taskIndex, err := support.ParseNaturalIndex(config.TaskIndex, 0, nt-1)
	if err != nil {
		return base{}, err
	}

	nu := uint(len(taskIndex))

	deviation := make([]float64, nu)
	for i, tid := range taskIndex {
		deviation[i] = config.Deviation * reference[tid]
	}

	return base{
		taskIndex: taskIndex,
		reference: reference,
		deviation: deviation,

		nt: nt,
		nu: nu,
		nz: nu,
	}, nil
}

func (self *base) Len() int {
	return int(self.nz)
}

func (self *base) String() string {
	return fmt.Sprintf(`{"variables": %d, "parameters": %d}`, self.nz, self.nu)
}
