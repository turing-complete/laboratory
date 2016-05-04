package uncertainty

import (
	"errors"
	"math"

	"github.com/ready-steady/infinity"
	"github.com/ready-steady/linear/matrix"
	"github.com/ready-steady/probability/distribution"
	"github.com/ready-steady/statistics/correlation"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/support"
	"github.com/turing-complete/laboratory/src/internal/system"

	icorrelation "github.com/turing-complete/laboratory/src/internal/correlation"
	idistribution "github.com/turing-complete/laboratory/src/internal/distribution"
)

var (
	epsilon  = math.Nextafter(1.0, 2.0) - 1.0
	gaussian = distribution.NewGaussian(0.0, 1.0)
)

type base struct {
	tasks []uint
	lower []float64
	upper []float64

	nt uint
	nu uint
	nz uint

	copula    *copula
	marginals []distribution.Continuous
}

// x(z) = F^(-1)(u(z))
// u(z) = Φ(C * Φ^(-1)(z))
//
// z(x) = Φ(D * Φ^(-1)(u(x)))
// u(x) = F(x)
//
// f(u) = exp(-0.5 * Φ^(-1)(u)^T * (R^(-1) - I) * Φ^(-1)(u)) / sqrt(det(R))
//      = exp(-0.5 * Φ^(-1)(u)^T * P * Φ^(-1)(u)) / N
//
// f(x) = prod(f(x)) * f(F(x))
type copula struct {
	R []float64
	C []float64
	D []float64
	P []float64
	N float64
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

	copula, err := correlate(system, config, tasks)
	if err != nil {
		return nil, err
	}

	nz := uint(len(copula.C)) / nu

	marginal, err := idistribution.Parse(config.Distribution)
	if err != nil {
		return nil, err
	}

	marginals := make([]distribution.Continuous, nu)
	for i := uint(0); i < nu; i++ {
		marginals[i] = marginal
	}

	return &base{
		tasks: tasks,
		lower: lower,
		upper: upper,

		nt: nt,
		nu: nu,
		nz: nz,

		copula:    copula,
		marginals: marginals,
	}, nil
}

func (self *base) Mapping() (uint, uint) {
	return self.nz, self.nt
}

func (self *base) Evaluate(ω []float64) float64 {
	nu, nz := self.nu, self.nz
	lower, upper := self.lower, self.upper

	if nu != nz {
		panic("model-order reduction is not supported")
	}

	amplitude, u := 1.0, make([]float64, nu)
	for i, tid := range self.tasks {
		ω := (ω[tid] - lower[tid]) / (upper[tid] - lower[tid])
		u[i] = gaussian.Invert(self.marginals[i].Cumulate(ω))
		amplitude *= self.marginals[i].Weigh(ω)
	}

	exponent := -0.5 * infinity.Quadratic(self.copula.P, u, nu)

	return amplitude * math.Exp(exponent) / self.copula.N
}

func (self *base) Forward(ω []float64) []float64 {
	nu, nz := self.nu, self.nz
	lower, upper := self.lower, self.upper

	z := make([]float64, nz)
	u := make([]float64, nu)

	// Dependent desired to dependent uniform
	for i, tid := range self.tasks {
		ω := (ω[tid] - lower[tid]) / (upper[tid] - lower[tid])
		u[i] = self.marginals[i].Cumulate(ω)
	}

	// Dependent uniform to dependent Gaussian
	for i := range u {
		u[i] = gaussian.Invert(u[i])
	}

	// Dependent Gaussian to independent Gaussian
	n := infinity.Linear(self.copula.D, u, nz, nu)

	// Independent Gaussian to independent uniform
	for i := range n {
		z[i] = gaussian.Cumulate(n[i])
	}

	return z
}

func (self *base) Backward(z []float64) []float64 {
	nu, nz := self.nu, self.nz
	lower, upper := self.lower, self.upper

	ω := append([]float64(nil), lower...)
	n := make([]float64, nz)

	// Independent uniform to independent Gaussian
	for i := range n {
		n[i] = gaussian.Invert(z[i])
	}

	// Independent Gaussian to dependent Gaussian
	u := infinity.Linear(self.copula.C, n, nu, nz)

	// Dependent Gaussian to dependent uniform
	for i := range u {
		u[i] = gaussian.Cumulate(u[i])
	}

	// Dependent uniform to dependent desired
	for i, tid := range self.tasks {
		ω[tid] += (upper[tid] - lower[tid]) * self.marginals[i].Invert(u[i])
	}

	return ω
}

func correlate(system *system.System, config *config.Uncertainty,
	tasks []uint) (*copula, error) {

	ε := math.Sqrt(epsilon)

	nu := uint(len(tasks))

	if config.Correlation == 0.0 {
		return &copula{
			R: matrix.Identity(nu),
			C: matrix.Identity(nu),
			D: matrix.Identity(nu),
			P: make([]float64, nu*nu),
			N: 1.0,
		}, nil
	}
	if config.Correlation < 0.0 {
		return nil, errors.New("the correlation length should be nonnegative")
	}
	if config.Variance <= 0.0 {
		return nil, errors.New("the variance threshold should be positive")
	}

	R := icorrelation.Compute(system.Application, tasks, config.Correlation)

	C, D, U, Λ, err := correlation.Decompose(R, nu, config.Variance, ε)
	if err != nil {
		return nil, err
	}

	detR := 1.0
	for _, λ := range Λ {
		if λ <= 0.0 {
			return nil, errors.New("the corelation matrix is invalid or singular")
		}
		detR *= λ
	}

	P, err := invert(U, Λ, nu)
	if err != nil {
		return nil, err
	}
	for i := uint(0); i < nu; i++ {
		P[i*nu+i] -= 1.0
	}

	return &copula{
		R: R,
		C: C,
		D: D,
		P: P,
		N: math.Sqrt(detR),
	}, nil
}
