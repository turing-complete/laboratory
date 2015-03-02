package internal

import (
	"errors"
	"fmt"

	"github.com/ready-steady/linear/matrix"
	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/gaussian"
	"github.com/ready-steady/simulation/system"
	"github.com/ready-steady/simulation/time"
	"github.com/ready-steady/statistics/correlation"

	acorrelation "../../pkg/correlation"
	aprobability "../../pkg/probability"
)

var standardGaussian = gaussian.New(0, 1)

type Problem struct {
	Config Config

	platform    *system.Platform
	application *system.Application

	nc uint
	nt uint
	nu uint
	nz uint

	marginals  []probability.Inverter
	multiplier []float64

	time     *time.List
	schedule *time.Schedule
}

func (p *Problem) String() string {
	return fmt.Sprintf("Problem{cores: %d, tasks: %d, dvars: %d, ivars: %d}",
		p.nc, p.nt, p.nu, p.nz)
}

func NewProblem(config Config) (*Problem, error) {
	p := &Problem{Config: config}
	c := &p.Config

	if c.Probability.MaxDelay < 0 || 1 <= c.Probability.MaxDelay {
		return nil, errors.New("the delay rate is invalid")
	}
	if c.Probability.CorrLength <= 0 {
		return nil, errors.New("the correlation length is invalid")
	}
	if c.Probability.VarThreshold <= 0 {
		return nil, errors.New("the variance-reduction threshold is invalid")
	}

	platform, application, err := system.Load(c.TGFF)
	if err != nil {
		return nil, err
	}

	p.platform = platform
	p.application = application

	p.nc = uint(len(platform.Cores))
	p.nt = uint(len(application.Tasks))

	if len(c.CoreIndex) == 0 {
		c.CoreIndex = make([]uint, p.nc)
		for i := uint(0); i < p.nc; i++ {
			c.CoreIndex[i] = i
		}
	}
	if len(c.TaskIndex) == 0 {
		c.TaskIndex = make([]uint, p.nt)
		for i := uint(0); i < p.nt; i++ {
			c.TaskIndex[i] = i
		}
	}

	p.nu = uint(len(c.TaskIndex))

	p.time = time.NewList(platform, application)
	p.schedule = p.time.Compute(system.NewProfile(platform, application).Mobility)

	C := acorrelation.Compute(application, c.TaskIndex, c.Probability.CorrLength)
	p.multiplier, p.nz, err = correlation.Decompose(C, p.nu, c.Probability.VarThreshold)
	if err != nil {
		return nil, err
	}

	p.marginals = make([]probability.Inverter, p.nu)
	marginalizer := aprobability.ParseInverter(c.Probability.Marginal)
	if marginalizer == nil {
		return nil, errors.New("invalid marginal distributions")
	}
	for i, tid := range c.TaskIndex {
		duration := platform.Cores[p.schedule.Mapping[tid]].Time[application.Tasks[tid].Type]
		p.marginals[i] = marginalizer(0, c.Probability.MaxDelay*duration)
	}

	return p, nil
}

func (p *Problem) transform(node []float64) []float64 {
	const (
		offset = 1e-8
	)

	z := make([]float64, p.nz)
	u := make([]float64, p.nu)
	d := make([]float64, p.nt)

	// Independent uniform to independent Gaussian
	for i := range z {
		switch node[i] {
		case 0:
			z[i] = standardGaussian.InvCDF(0 + offset)
		case 1:
			z[i] = standardGaussian.InvCDF(1 - offset)
		default:
			z[i] = standardGaussian.InvCDF(node[i])
		}
	}

	// Independent Gaussian to dependent Gaussian
	matrix.Multiply(p.multiplier, z, u, p.nu, p.nz, 1)

	// Dependent Gaussian to dependent uniform to dependent target
	for i, tid := range p.Config.TaskIndex {
		d[tid] = p.marginals[i].InvCDF(standardGaussian.CDF(u[i]))
	}

	return d
}
