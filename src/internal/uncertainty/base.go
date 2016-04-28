package uncertainty

import (
	"errors"
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
	standardGaussian = probability.NewGaussian(0.0, 1.0)
	nInfinity        = math.Inf(-1.0)
	pInfinity        = math.Inf(1.0)
)

type base struct {
	tasks []uint
	lower []float64
	upper []float64

	nt uint
	nu uint
	nz uint

	correlator   []float64
	decorrelator []float64
	marginals    []probability.Distribution
}

func newBase(system *system.System, reference []float64,
	config *config.Uncertainty) (*base, error) {

	nt := uint(len(reference))

	tasks, err := support.ParseNaturalIndex(config.Tasks, 0, nt-1)
	if err != nil {
		return nil, err
	}

	nu := uint(len(tasks))

	lower := make([]float64, nt)
	upper := make([]float64, nt)

	copy(lower, reference)
	copy(upper, reference)

	for _, tid := range tasks {
		lower[tid] -= config.Deviation * reference[tid]
		upper[tid] += config.Deviation * reference[tid]
	}

	if nu == 0 {
		return &base{
			tasks: tasks,
			lower: lower,
			upper: upper,

			nt: nt,
		}, nil
	}

	correlator, decorrelator, err := correlate(system, config, tasks)
	if err != nil {
		return nil, err
	}

	marginalizer, err := distribution.Parse(config.Distribution)
	if err != nil {
		return nil, err
	}

	marginals := make([]probability.Distribution, nu)
	for i, tid := range tasks {
		marginals[i] = marginalizer(lower[tid], upper[tid])
	}

	return &base{
		tasks: tasks,
		lower: lower,
		upper: upper,

		nt: nt,
		nu: nu,
		nz: uint(len(correlator)) / nu,

		correlator:   correlator,
		decorrelator: decorrelator,
		marginals:    marginals,
	}, nil
}

func (self *base) Mapping() (uint, uint) {
	return self.nz, self.nt
}

func (self *base) Forward(ω []float64) []float64 {
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

func (self *base) Inverse(z []float64) []float64 {
	nu, nz := self.nu, self.nz

	ω := append([]float64(nil), self.lower...)
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

func correlate(system *system.System, config *config.Uncertainty,
	tasks []uint) ([]float64, []float64, error) {

	if config.Correlation == 0.0 {
		identity := matrix.Identity(uint(len(tasks)))
		return identity, append([]float64(nil), identity...), nil
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
