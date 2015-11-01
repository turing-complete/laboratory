package uncertainty

import (
	"errors"
	"fmt"
	"math"

	"github.com/ready-steady/linear/matrix"
	"github.com/ready-steady/probability"
	"github.com/ready-steady/statistics/correlation"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/distribution"
	"github.com/turing-complete/laboratory/src/internal/support"
	"github.com/turing-complete/laboratory/src/internal/system"

	icorrelation "github.com/turing-complete/laboratory/src/internal/correlation"
)

var (
	standardGaussian = probability.NewGaussian(0, 1)
	nInfinity        = math.Inf(-1)
	pInfinity        = math.Inf(1)
)

type marginal struct {
	taskIndex  []uint
	correlator []float64
	marginals  []probability.Inverter
	reference  []float64

	nt uint
	nu uint
	nz uint
}

func newMarginal(s *system.System, c *config.Uncertainty) (*marginal, error) {
	nt := uint(s.Application.Len())

	taskIndex, err := support.ParseNaturalIndex(c.TaskIndex, 0, nt-1)
	if err != nil {
		return nil, err
	}

	correlator, err := correlate(s, c, taskIndex)
	if err != nil {
		return nil, err
	}

	nu := uint(len(taskIndex))
	nz := uint(len(correlator)) / nu

	marginalizer, err := distribution.ParseInverter(c.Distribution)
	if err != nil {
		return nil, err
	}

	reference := s.ReferenceTime()
	marginals := make([]probability.Inverter, nu)
	for i, tid := range taskIndex {
		marginals[i] = marginalizer(0, c.MaxDeviation*reference[tid])
	}

	return &marginal{
		taskIndex:  taskIndex,
		correlator: correlator,
		marginals:  marginals,
		reference:  reference,

		nt: nt,
		nu: nu,
		nz: nz,
	}, nil
}

func (m *marginal) Transform(z []float64) []float64 {
	nt, nu, nz := m.nt, m.nu, m.nz

	n := make([]float64, nz)
	u := make([]float64, nu)

	// Independent uniform to independent Gaussian
	for i := range n {
		n[i] = standardGaussian.InvCDF(z[i])
	}

	// Independent Gaussian to dependent Gaussian
	multiply(m.correlator, n, u, nu, nz)

	// Dependent Gaussian to dependent uniform
	for i := range u {
		u[i] = standardGaussian.CDF(u[i])
	}

	// Dependent uniform to dependent desired
	duration := make([]float64, nt)
	copy(duration, m.reference)
	for i, tid := range m.taskIndex {
		duration[tid] += m.marginals[i].InvCDF(standardGaussian.CDF(u[i]))
	}

	return duration
}

func (m *marginal) Len() int {
	return int(m.nz)
}

func (m *marginal) String() string {
	return fmt.Sprintf(`{"parameters": %d, "variables": %d}`, m.nu, m.nz)
}

func correlate(s *system.System, c *config.Uncertainty, taskIndex []uint) ([]float64, error) {
	if c.CorrLength < 0 {
		return nil, errors.New("the correlation length should be nonnegative")
	}
	if c.VarThreshold <= 0 {
		return nil, errors.New("the variance-reduction threshold should be positive")
	}

	if c.CorrLength == 0 {
		return matrix.Identity(uint(len(taskIndex))), nil
	}

	C := icorrelation.Compute(s.Application, taskIndex, c.CorrLength)
	correlator, _, err := correlation.Decompose(C, uint(len(taskIndex)), c.VarThreshold)
	if err != nil {
		return nil, err
	}

	return correlator, nil
}

func multiply(A, x, y []float64, m, n uint) {
	infinite, z := false, make([]float64, n)

	for i := range x {
		switch x[i] {
		case nInfinity:
			infinite, z[i] = true, -1
		case pInfinity:
			infinite, z[i] = true, 1
		}
	}

	if !infinite {
		matrix.Multiply(A, x, y, m, n, 1)
		return
	}

	for i := uint(0); i < m; i++ {
		Σ1, Σ2 := 0.0, 0.0
		for j := uint(0); j < n; j++ {
			a := A[j*m+i]
			if a == 0 {
				continue
			}
			if z[j] == 0 {
				Σ1 += a * x[j]
			} else {
				Σ2 += a * z[j]
			}
		}
		if Σ2 < 0 {
			y[i] = nInfinity
		} else if Σ2 > 0 {
			y[i] = pInfinity
		} else {
			y[i] = Σ1
		}
	}
}
