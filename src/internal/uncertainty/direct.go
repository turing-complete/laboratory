package uncertainty

import (
	"fmt"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/support"
	"github.com/turing-complete/laboratory/src/internal/system"
)

type direct struct {
	taskIndex []uint
	reference []float64
	deviation []float64

	nt uint
	nu uint
}

func newDirect(system *system.System, config *config.Uncertainty) (*direct, error) {
	nt := uint(system.Application.Len())

	taskIndex, err := support.ParseNaturalIndex(config.TaskIndex, 0, nt-1)
	if err != nil {
		return nil, err
	}

	nu := uint(len(taskIndex))

	reference := system.ReferenceTime()
	deviation := make([]float64, nu)
	for i, tid := range taskIndex {
		deviation[i] = config.MaxDeviation * reference[tid]
	}

	return &direct{
		taskIndex: taskIndex,
		reference: reference,
		deviation: deviation,

		nt: nt,
		nu: nu,
	}, nil
}

func (m *direct) Transform(z []float64) []float64 {
	duration := make([]float64, m.nt)
	copy(duration, m.reference)
	for i, tid := range m.taskIndex {
		duration[tid] += z[i] * m.deviation[i]
	}
	return duration
}

func (m *direct) Len() int {
	return int(m.nu)
}

func (m *direct) String() string {
	return fmt.Sprintf(`{"parameters": %d}`, m.nu)
}
