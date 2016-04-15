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
	standardGaussian = probability.NewGaussian(0, 1)
	nInfinity        = math.Inf(-1)
	pInfinity        = math.Inf(1)
)

type aleatory struct {
	base
	correlator []float64
	marginals  []probability.Inverter

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

	correlator, err := correlate(system, config, base.tasks)
	if err != nil {
		return nil, err
	}

	marginalizer, err := distribution.ParseInverter(config.Distribution)
	if err != nil {
		return nil, err
	}

	marginals := make([]probability.Inverter, base.nu)
	for i, tid := range base.tasks {
		marginals[i] = marginalizer(0, base.upper[tid]-base.lower[tid])
	}

	return &aleatory{
		base:       base,
		correlator: correlator,
		marginals:  marginals,

		nz: uint(len(correlator)) / base.nu,
	}, nil
}

func (self *aleatory) Mapping() (uint, uint) {
	return self.nz, self.nt
}

func (_ *aleatory) Forward(_ []float64) []float64 {
	return nil
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
		n[i] = standardGaussian.InvCDF(z[i])
	}

	// Independent Gaussian to dependent Gaussian
	multiply(self.correlator, n, u, nu, nz)

	// Dependent Gaussian to dependent uniform
	for i := range u {
		u[i] = standardGaussian.CDF(u[i])
	}

	// Dependent uniform to dependent desired
	for i, tid := range self.tasks {
		ω[tid] += self.marginals[i].InvCDF(standardGaussian.CDF(u[i]))
	}

	return ω
}

func correlate(system *system.System, config *config.Parameter, tasks []uint) ([]float64, error) {
	if config.Correlation < 0 {
		return nil, errors.New("the correlation length should be nonnegative")
	}
	if config.Variance <= 0 {
		return nil, errors.New("the variance threshold should be positive")
	}

	if config.Correlation == 0 {
		return matrix.Identity(uint(len(tasks))), nil
	}

	C := icorrelation.Compute(system.Application, tasks, config.Correlation)
	correlator, _, err := correlation.Decompose(C, uint(len(tasks)), config.Variance)
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