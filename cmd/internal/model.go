package internal

import (
	"errors"
	"fmt"
	"sort"

	"github.com/ready-steady/probability"
	"github.com/ready-steady/staircase"
	"github.com/ready-steady/statistics/correlation"
	"github.com/simulated-reality/laboratory/internal/config"
	"github.com/simulated-reality/laboratory/internal/support"

	acorrelation "github.com/simulated-reality/laboratory/internal/correlation"
)

var (
	standardGaussian = probability.NewGaussian(0, 1)
)

type model struct {
	taskIndex  []uint
	correlator []float64
	modes      []mode

	nt uint
	nu uint
	nz uint
}

type mode *staircase.Staircase

func newModel(c *config.Probability, s *system) (*model, error) {
	nt := s.nt

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
	return fmt.Sprintf(`{"parameters": %d, "variables": %d}`, m.nu, m.nz)
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
	support.Combine(m.correlator, n, u, nu, nz)

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

func computeCorrelator(c *config.Probability, s *system, taskIndex []uint) ([]float64, error) {
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

func computeModes(c *config.Probability, count uint) ([]mode, error) {
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
