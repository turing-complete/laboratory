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
	reference []float64
}

func newMarginal(s *system.System, c *config.Uncertainty) (*Marginal, error) {
	base, err := newBase(s, c)
	if err != nil {
		return nil, err
	}

	marginalizer, err := distribution.ParseInverter(c.Marginal)
	if err != nil {
		return nil, err
	}

	reference := s.ReferenceTime()
	marginals := make([]probability.Inverter, base.nu)
	for i, tid := range base.taskIndex {
		marginals[i] = marginalizer(0, c.MaxDelay*reference[tid])
	}

	return &Marginal{
		base:      *base,
		marginals: marginals,
		reference: reference,
	}, nil
}

func (m *Marginal) Transform(z []float64) []float64 {
	u := m.base.Transform(z)

	duration := make([]float64, m.nt)
	copy(duration, m.reference)
	for i, tid := range m.taskIndex {
		duration[tid] += m.marginals[i].InvCDF(standardGaussian.CDF(u[i]))
	}

	return duration
}
