package uncertainty

import (
	"errors"
	"fmt"
	"sort"

	"github.com/ready-steady/probability"
	"github.com/ready-steady/staircase"
	"github.com/simulated-reality/laboratory/src/internal/config"
	"github.com/simulated-reality/laboratory/src/internal/support"
	"github.com/simulated-reality/laboratory/src/internal/system"
)

var (
	standardGaussian = probability.NewGaussian(0, 1)
)

type Modal struct {
	taskIndex  []uint
	correlator []float64
	modes      []mode

	nt uint
	nu uint
	nz uint
}

type mode *staircase.Staircase

func NewModal(c *config.Uncertainty, s *system.System) (*Modal, error) {
	nt := uint(s.Application.Len())

	taskIndex, err := support.ParseNaturalIndex(c.TaskIndex, 0, nt-1)
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

	model := &Modal{
		taskIndex:  taskIndex,
		correlator: correlator,
		modes:      modes,

		nt: nt,
		nu: nu,
		nz: nz,
	}

	return model, nil
}

func (m *Modal) Len() int {
	return int(m.nz)
}

func (m *Modal) String() string {
	return fmt.Sprintf(`{"parameters": %d, "variables": %d}`, m.nu, m.nz)
}

func (m *Modal) Transform(z []float64) []float64 {
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
		modes[tid] = (*staircase.Staircase)(m.modes[i]).Evaluate(u[i])
	}

	return modes
}

func computeModes(c *config.Uncertainty, count uint) ([]mode, error) {
	if c.Modes == 0 {
		return nil, errors.New("the number of modes should be positive")
	}
	if c.MinOffset <= -1 || c.MaxOffset <= -1 {
		return nil, errors.New("the offsets should be greater than -1")
	}
	if c.Transition <= 0 || c.Transition > 0.5 {
		return nil, errors.New("the transition parameter should be in (0, 0.5]")
	}

	generator := probability.NewGenerator(support.NewSeed(c.Seed))
	uniform := probability.NewUniform(0, 1)

	result := make([]mode, count)
	for i := range result {
		Σ := 0.0
		values := make([]float64, c.Modes)
		probabilities := make([]float64, c.Modes)
		for j := range values {
			values[j] = uniform.Sample(generator)
			probabilities[j] = uniform.Sample(generator)
			Σ += probabilities[j]
		}
		for j := range values {
			values[j] = c.MinOffset + (c.MaxOffset-c.MinOffset)*values[j]
			probabilities[j] /= Σ
		}
		sort.Float64s(values)

		result[i] = staircase.New(probabilities, values, c.Transition)
	}

	return result, nil
}
