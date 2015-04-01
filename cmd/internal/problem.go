package internal

import (
	"errors"
	"fmt"

	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/gaussian"
	"github.com/ready-steady/statistics/correlation"

	acorrelation "../../pkg/correlation"
	aprobability "../../pkg/probability"
)

var (
	standardGaussian = gaussian.New(0, 1)
)

type Problem struct {
	Config *Config
	system *system

	taskIndex  []uint
	multiplier []float64
	marginals  []probability.Inverter

	nu uint
	nz uint
}

func (p *Problem) String() string {
	return fmt.Sprintf("Problem{cores: %d, tasks: %d, dvars: %d, ivars: %d}",
		p.system.nc, p.system.nt, p.nu, p.nz)
}

func NewProblem(config *Config) (*Problem, error) {
	system, err := newSystem(&config.System)
	if err != nil {
		return nil, err
	}

	c := &config.Probability
	if c.MaxDelay < 0 {
		return nil, errors.New("the maximal delay should be nonnegative")
	}
	if c.CorrLength < 0 {
		return nil, errors.New("the correlation length should be nonnegative")
	}
	if c.VarThreshold <= 0 {
		return nil, errors.New("the variance-reduction threshold should be positive")
	}

	taskIndex, err := parseNaturalIndex(c.TaskIndex, 0, system.nt-1)
	if err != nil {
		return nil, err
	}

	nu := uint(len(taskIndex))

	C := acorrelation.Compute(system.application, taskIndex, c.CorrLength)
	multiplier, nz, err := correlation.Decompose(C, nu, c.VarThreshold)
	if err != nil {
		return nil, err
	}

	marginalizer, err := aprobability.ParseInverter(c.Marginal)
	if err != nil {
		return nil, err
	}
	duration := system.computeTime(system.schedule)
	marginals := make([]probability.Inverter, nu)
	for i, j := range taskIndex {
		marginals[i] = marginalizer(0, c.MaxDelay*duration[j])
	}

	problem := &Problem{
		Config: config,
		system: system,

		taskIndex:  taskIndex,
		multiplier: multiplier,
		marginals:  marginals,

		nu: nu,
		nz: nz,
	}

	return problem, nil
}

func (p *Problem) transform(z []float64) []float64 {
	nt, nu, nz := p.system.nt, p.nu, p.nz

	n := make([]float64, nz)
	u := make([]float64, nu)
	d := make([]float64, nt)

	// Independent uniform to independent Gaussian
	for i := range n {
		n[i] = standardGaussian.InvCDF(z[i])
	}

	// Independent Gaussian to dependent Gaussian
	combine(p.multiplier, n, u, nu, nz)

	// Dependent Gaussian to dependent uniform to dependent target
	for i, j := range p.taskIndex {
		d[j] = p.marginals[i].InvCDF(standardGaussian.CDF(u[i]))
	}

	return d
}
