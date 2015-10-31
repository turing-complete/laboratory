package uncertainty

import (
	"errors"
	"sort"

	"github.com/ready-steady/probability"
	"github.com/ready-steady/staircase"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/support"
	"github.com/turing-complete/laboratory/src/internal/system"
)

type Modal struct {
	base
	modes     []mode
	reference []float64
}

type mode *staircase.Staircase

func NewModal(s *system.System, c *config.Uncertainty) (*Modal, error) {
	base, err := newBase(c, s)
	if err != nil {
		return nil, err
	}

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

	modes := make([]mode, base.nu)
	for i := range modes {
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

		modes[i] = staircase.New(probabilities, values, c.Transition)
	}

	return &Modal{
		base:      *base,
		modes:     modes,
		reference: s.ReferenceTime(),
	}, nil
}

func (m *Modal) Transform(z []float64) []float64 {
	u := m.base.Transform(z)

	duration := make([]float64, m.nt)
	copy(duration, m.reference)
	for i, tid := range m.taskIndex {
		duration[tid] *= 1.0 + (*staircase.Staircase)(m.modes[i]).Evaluate(u[i])
	}

	return duration
}
