package uncertainty

import (
	"errors"
	"math"

	"github.com/ready-steady/linear/matrix"
	"github.com/ready-steady/probability"
	"github.com/ready-steady/statistics/correlation"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/distribution"
	"github.com/turing-complete/laboratory/src/internal/system"

	icorrelation "github.com/turing-complete/laboratory/src/internal/correlation"
)

var (
	standardGaussian = probability.NewGaussian(0.0, 1.0)
	nInfinity        = math.Inf(-1.0)
	pInfinity        = math.Inf(1.0)
)

type aleatory struct {
	base

	correlator   []float64
	decorrelator []float64
	marginals    []probability.Distribution

	nz uint
}

func newAleatory(system *system.System, reference []float64,
	config *config.Parameter) (*aleatory, error) {

	base, err := newBase(reference, config)
	if err != nil {
		return nil, err
	}

	if base.nu == 0 {
		return &aleatory{base: base}, nil
	}

	correlator, decorrelator, err := correlate(system, config, base.tasks)
	if err != nil {
		return nil, err
	}

	marginalizer, err := distribution.Parse(config.Distribution)
	if err != nil {
		return nil, err
	}

	marginals := make([]probability.Distribution, base.nu)
	for i, tid := range base.tasks {
		marginals[i] = marginalizer(base.lower[tid], base.upper[tid])
	}

	return &aleatory{
		base: base,

		correlator:   correlator,
		decorrelator: decorrelator,
		marginals:    marginals,

		nz: uint(len(correlator)) / base.nu,
	}, nil
}

func (self *aleatory) Mapping() (uint, uint) {
	return self.nz, self.nt
}

func (self *aleatory) Forward(ω []float64) []float64 {
	nu, nz := self.nu, self.nz

	z := make([]float64, nz)
	if nz == 0 {
		return z
	}

	u := make([]float64, nu)
	n := make([]float64, nz)

	// Dependent desired to dependent uniform
	for i, tid := range self.tasks {
		u[i] = self.marginals[i].Cumulate(ω[tid])
	}

	// Dependent uniform to dependent Gaussian
	for i := range u {
		u[i] = standardGaussian.Decumulate(u[i])
	}

	// Dependent Gaussian to independent Gaussian
	multiply(self.decorrelator, u, n, nz, nu)

	// Independent Gaussian to independent uniform
	for i := range n {
		z[i] = standardGaussian.Cumulate(n[i])
	}

	return z
}

func (self *aleatory) Inverse(z []float64) []float64 {
	nt, nu, nz := self.nt, self.nu, self.nz

	ω := make([]float64, nt)
	copy(ω, self.lower)
	if nu == 0 {
		return ω
	}

	n := make([]float64, nz)
	u := make([]float64, nu)

	// Independent uniform to independent Gaussian
	for i := range n {
		n[i] = standardGaussian.Decumulate(z[i])
	}

	// Independent Gaussian to dependent Gaussian
	multiply(self.correlator, n, u, nu, nz)

	// Dependent Gaussian to dependent uniform
	for i := range u {
		u[i] = standardGaussian.Cumulate(u[i])
	}

	// Dependent uniform to dependent desired
	for i, tid := range self.tasks {
		ω[tid] = self.marginals[i].Decumulate(u[i])
	}

	return ω
}

func correlate(system *system.System, config *config.Parameter,
	tasks []uint) ([]float64, []float64, error) {

	if config.Correlation == 0.0 {
		identity := matrix.Identity(uint(len(tasks)))
		return identity, append(([]float64)(nil), identity...), nil
	}
	if config.Correlation < 0.0 {
		return nil, nil, errors.New("the correlation length should be nonnegative")
	}
	if config.Variance <= 0.0 {
		return nil, nil, errors.New("the variance threshold should be positive")
	}

	C := icorrelation.Compute(system.Application, tasks, config.Correlation)
	correlator, decorrelator, _, err := correlation.Decompose(C, uint(len(tasks)), config.Variance)
	if err != nil {
		return nil, nil, err
	}

	return correlator, decorrelator, nil
}

func multiply(A, x, y []float64, m, n uint) {
	infinite, z := false, make([]float64, n)

	for i := range x {
		switch x[i] {
		case nInfinity:
			infinite, z[i] = true, -1.0
		case pInfinity:
			infinite, z[i] = true, 1.0
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
			if a == 0.0 {
				continue
			}
			if z[j] == 0.0 {
				Σ1 += a * x[j]
			} else {
				Σ2 += a * z[j]
			}
		}
		if Σ2 < 0.0 {
			y[i] = nInfinity
		} else if Σ2 > 0.0 {
			y[i] = pInfinity
		} else {
			y[i] = Σ1
		}
	}
}
