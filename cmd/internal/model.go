package internal

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ready-steady/probability/gaussian"
	"github.com/ready-steady/probability/generator"
	"github.com/ready-steady/probability/uniform"
	"github.com/ready-steady/statistics/correlation"

	acorrelation "../../pkg/correlation"
)

var (
	standardGaussian = gaussian.New(0, 1)
)

type model struct {
	taskIndex  []uint
	correlator []float64
	modes      [][]mode

	nt uint
	nu uint
	nz uint
}

type mode struct {
	value float64
	point float64
}

func newModel(c *ProbabilityConfig, s *system) (*model, error) {
	nt := s.nt

	taskIndex, err := parseNaturalIndex(c.TaskIndex, 0, nt-1)
	if err != nil {
		return nil, err
	}

	correlator, err := computeCorrelator(c, s, taskIndex)
	if err != nil {
		return nil, err
	}

	nu := uint(len(taskIndex))
	nz := uint(len(correlator)) / nu

	modes, err := computeModes(c, nu)
	if err != nil {
		return nil, err
	}

	model := &model{
		taskIndex:  taskIndex,
		correlator: correlator,
		modes:      modes,

		nt: nt,
		nu: nu,
		nz: nz,
	}

	return model, nil
}

func (m *model) String() string {
	var buffer bytes.Buffer

	put := func(format string, arguments ...interface{}) {
		buffer.WriteString(fmt.Sprintf(format, arguments...))
	}

	put("[")
	for i, tid := range m.taskIndex {
		if i > 0 {
			put(", ")
		}
		put(`{"id": %d, "modes": [`, tid)
		for j := range m.modes[i] {
			if j > 0 {
				put(", ")
			}
			put("[%.2f, %.2f]", m.modes[i][j].value, m.modes[i][j].point)
		}
		put("]}")
	}
	put("]")

	return fmt.Sprintf(`{"parameters": %d, "variables": %d, "tasks": %s}`,
		m.nu, m.nz, buffer.String())
}

func (m *model) transform(z []float64) []float64 {
	nt, nu, nz := m.nt, m.nu, m.nz

	n := make([]float64, nz)
	u := make([]float64, nu)

	// Independent uniform to independent Gaussian
	for i := range n {
		n[i] = standardGaussian.InvCDF(z[i])
	}

	// Independent Gaussian to dependent Gaussian
	combine(m.correlator, n, u, nu, nz)

	// Dependent Gaussian to dependent uniform
	for i := range u {
		u[i] = standardGaussian.CDF(u[i])
	}

	modes := make([]float64, nt)
	for i, tid := range m.taskIndex {
		for j := range m.modes[i] {
			if u[i] <= m.modes[i][j].point {
				modes[tid] = m.modes[i][j].value
				break
			}
		}
	}

	return modes
}

func computeCorrelator(c *ProbabilityConfig, s *system, taskIndex []uint) ([]float64, error) {
	if c.CorrLength < 0 {
		return nil, errors.New("the correlation length should be nonnegative")
	}
	if c.VarThreshold <= 0 {
		return nil, errors.New("the variance-reduction threshold should be positive")
	}

	C := acorrelation.Compute(s.application, taskIndex, c.CorrLength)
	correlator, _, err := correlation.Decompose(C, uint(len(taskIndex)), c.VarThreshold)
	if err != nil {
		return nil, err
	}

	return correlator, nil
}

func computeModes(c *ProbabilityConfig, count uint) ([][]mode, error) {
	if c.MaxModes == 0 {
		return nil, errors.New("the number of modes should be positive")
	}
	if c.MinScale <= -1 || c.MaxScale <= -1 {
		return nil, errors.New("the scaling factors should be greater than -1")
	}

	generator := generator.New(NewSeed(c.Seed))
	uniform := uniform.New(0, 1)

	result := make([][]mode, count)
	for i := range result {
		count := uint(uniform.Sample(generator)*float64(c.MaxModes)) + 1

		result[i] = make([]mode, count)
		for j := range result[i] {
			result[i][j].value = uniform.Sample(generator)
			result[i][j].point = uniform.Sample(generator)
			if j > 0 {
				result[i][j].point += result[i][j-1].point
			}
		}
		for j := range result[i] {
			result[i][j].value *= c.MaxScale - c.MinScale
			result[i][j].value += c.MinScale
			result[i][j].point /= result[i][count-1].point
		}
	}

	return result, nil
}
