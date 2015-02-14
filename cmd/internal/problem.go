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

	cc uint
	tc uint
	uc uint
	zc uint

	marginals  []probability.Inverter
	multiplier []float64

	time     *time.List
	schedule *time.Schedule
}

func (p *Problem) String() string {
	return fmt.Sprintf("Problem{cores: %d, tasks: %d, dvars: %d, ivars: %d}",
		p.cc, p.tc, p.uc, p.zc)
}

func NewProblem(config Config) (*Problem, error) {
	p := &Problem{Config: config}
	c := &p.Config

	if c.ProbModel.MaxDelay < 0 || 1 <= c.ProbModel.MaxDelay {
		return nil, errors.New("the delay rate is invalid")
	}
	if c.ProbModel.CorrLength <= 0 {
		return nil, errors.New("the correlation length is invalid")
	}
	if c.ProbModel.VarThreshold <= 0 {
		return nil, errors.New("the variance-reduction threshold is invalid")
	}

	platform, application, err := system.Load(c.TGFF)
	if err != nil {
		return nil, err
	}

	p.platform = platform
	p.application = application

	p.cc = uint(len(platform.Cores))
	p.tc = uint(len(application.Tasks))

	if len(c.CoreIndex) == 0 {
		c.CoreIndex = make([]uint, p.cc)
		for i := uint(0); i < uint(p.cc); i++ {
			c.CoreIndex[i] = i
		}
	}
	if len(c.TaskIndex) == 0 {
		c.TaskIndex = make([]uint, p.tc)
		for i := uint(0); i < uint(p.tc); i++ {
			c.TaskIndex[i] = i
		}
	}

	p.uc = uint(len(c.TaskIndex))

	p.time = time.NewList(platform, application)
	p.schedule = p.time.Compute(system.NewProfile(platform, application).Mobility)

	C := acorrelation.Compute(application, c.TaskIndex, c.ProbModel.CorrLength)
	p.multiplier, p.zc, err = correlation.Decompose(C, p.uc, c.ProbModel.VarThreshold)
	if err != nil {
		return nil, err
	}

	p.marginals = make([]probability.Inverter, p.uc)
	marginalizer := aprobability.ParseInverter(c.ProbModel.Marginal)
	if marginalizer == nil {
		return nil, errors.New("invalid marginal distributions")
	}
	for i, tid := range c.TaskIndex {
		duration := platform.Cores[p.schedule.Mapping[tid]].Time[application.Tasks[tid].Type]
		p.marginals[i] = marginalizer(0, c.ProbModel.MaxDelay*duration)
	}

	return p, nil
}

func (p *Problem) transform(node []float64) []float64 {
	const (
		offset = 1e-8
	)

	z := make([]float64, p.zc)
	u := make([]float64, p.uc)
	d := make([]float64, p.tc)

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
	matrix.Multiply(p.multiplier, z, u, p.uc, p.zc, 1)

	// Dependent Gaussian to dependent uniform to dependent target
	for i, tid := range p.Config.TaskIndex {
		d[tid] = p.marginals[i].InvCDF(standardGaussian.CDF(u[i]))
	}

	return d
}
