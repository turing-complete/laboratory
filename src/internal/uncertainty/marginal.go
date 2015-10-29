package uncertainty

import (
	"github.com/ready-steady/probability"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/distribution"
	"github.com/turing-complete/laboratory/src/internal/system"
)

type Marginal struct {
	base
	marginals []probability.Inverter
}

func NewMarginal(c *config.Uncertainty, s *system.System) (*Marginal, error) {
	base, err := newBase(c, s)
	if err != nil {
		return nil, err
	}

	marginalizer, err := distribution.ParseInverter(c.Marginal)
	if err != nil {
		return nil, err
	}

	reference := s.ReferenceTime()
	marginals := make([]probability.Inverter, base.nu)
	for i, j := range base.taskIndex {
		marginals[i] = marginalizer(0, c.MaxDelay*reference[j])
	}

	return &Marginal{base: *base, marginals: marginals}, nil
}

func (m *Marginal) Transform(z []float64) []float64 {
	u := m.base.Transform(z)

	duration := make([]float64, m.nt)
	for i, j := range m.taskIndex {
		duration[j] = m.marginals[i].InvCDF(standardGaussian.CDF(u[i]))
	}

	return duration
}
