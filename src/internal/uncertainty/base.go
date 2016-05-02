package uncertainty

import (
	"errors"
	"math"

	"github.com/ready-steady/linear/matrix"
	"github.com/ready-steady/probability/distribution"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/support"
	"github.com/turing-complete/laboratory/src/internal/system"

	scorrelation "github.com/ready-steady/statistics/correlation"
	icorrelation "github.com/turing-complete/laboratory/src/internal/correlation"
	idistribution "github.com/turing-complete/laboratory/src/internal/distribution"
)

var (
	epsilon          = math.Nextafter(1.0, 2.0) - 1.0
	standardGaussian = distribution.NewGaussian(0.0, 1.0)
)

type base struct {
	tasks []uint
	lower []float64
	upper []float64

	nt uint
	nu uint
	nz uint

	correlation *correlation
	marginals   []distribution.Continuous
}

type correlation struct {
	R []float64
	C []float64 // x = C * z
	D []float64 // z = D * x
	P []float64 // R^(-1) - I

	detR float64
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

	correlation, err := correlate(system, config, tasks)
	if err != nil {
		return nil, err
	}

	nz := uint(len(correlation.C)) / nu

	marginalizer, err := idistribution.Parse(config.Distribution)
	if err != nil {
		return nil, err
	}

	marginals := make([]distribution.Continuous, nu)
	for i, tid := range tasks {
		marginals[i] = marginalizer(lower[tid], upper[tid])
	}

	return &base{
		tasks: tasks,
		lower: lower,
		upper: upper,

		nt: nt,
		nu: nu,
		nz: nz,

		correlation: correlation,
		marginals:   marginals,
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
		u[i] = standardGaussian.Invert(u[i])
	}

	// Dependent Gaussian to independent Gaussian
	multiply(self.correlation.D, u, n, nz, nu)

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
		n[i] = standardGaussian.Invert(z[i])
	}

	// Independent Gaussian to dependent Gaussian
	multiply(self.correlation.C, n, u, nu, nz)

	// Dependent Gaussian to dependent uniform
	for i := range u {
		u[i] = standardGaussian.Cumulate(u[i])
	}

	// Dependent uniform to dependent desired
	for i, tid := range self.tasks {
		ω[tid] = self.marginals[i].Invert(u[i])
	}

	return ω
}

func (self *base) Weigh(ω []float64) float64 {
	nu, nz := self.nu, self.nz

	if nu != nz {
		panic("model-order reduction is not supported")
	}

	u := make([]float64, nu)
	n := make([]float64, nz)

	// Dependent desired to dependent uniform
	for i, tid := range self.tasks {
		u[i] = self.marginals[i].Cumulate(ω[tid])
	}

	// Dependent uniform to dependent Gaussian
	for i := range u {
		u[i] = standardGaussian.Invert(u[i])
	}

	multiply(self.correlation.P, u, n, nz, nu)

	amplitude := 1.0
	for i, tid := range self.tasks {
		amplitude *= self.marginals[i].Weigh(ω[tid])
	}
	amplitude = math.Abs(amplitude)

	exponent := 0.0
	for i := range n {
		exponent += u[i] * n[i]
	}
	exponent *= -0.5

	normalization := math.Sqrt(self.correlation.detR)

	return amplitude * math.Exp(exponent) / normalization
}

func correlate(system *system.System, config *config.Uncertainty,
	tasks []uint) (*correlation, error) {

	ε := math.Sqrt(epsilon)

	nu := uint(len(tasks))

	if config.Correlation == 0.0 {
		return &correlation{
			R: matrix.Identity(nu),
			C: matrix.Identity(nu),
			D: matrix.Identity(nu),
			P: make([]float64, nu*nu),

			detR: 1.0,
		}, nil
	}
	if config.Correlation < 0.0 {
		return nil, errors.New("the correlation length should be nonnegative")
	}
	if config.Variance <= 0.0 {
		return nil, errors.New("the variance threshold should be positive")
	}

	R := icorrelation.Compute(system.Application, tasks, config.Correlation)

	C, D, U, Λ, err := scorrelation.Decompose(R, nu, config.Variance, ε)
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

	P, err := inverse(U, Λ, nu)
	if err != nil {
		return nil, err
	}
	for i := uint(0); i < nu; i++ {
		P[i*nu+i] -= 1.0
	}

	return &correlation{
		R: R,
		C: C,
		D: D,
		P: P,

		detR: detR,
	}, nil
}
