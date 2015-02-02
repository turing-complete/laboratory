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
	config *Config

	platform    *system.Platform
	application *system.Application

	cc uint32
	tc uint32
	zc uint32

	marginals  []probability.Inverter
	multiplier []float64

	time     *time.List
	schedule *time.Schedule
}

func (p *Problem) String() string {
	return fmt.Sprintf("Problem{cores: %d, tasks: %d, variables: %d}",
		p.cc, p.tc, p.zc)
}

func newProblem(c *Config) (*Problem, error) {
	if c.ProbModel.MaxDelay < 0 || 1 <= c.ProbModel.MaxDelay {
		return nil, errors.New("the delay rate is invalid")
	}
	if c.ProbModel.CorrLength <= 0 {
		return nil, errors.New("the correlation length is invalid")
	}
	if c.ProbModel.VarThreshold <= 0 || 1 < c.ProbModel.VarThreshold {
		return nil, errors.New("the variance-reduction threshold is invalid")
	}

	var err error

	p := &Problem{config: c}

	platform, application, err := system.Load(c.TGFF)
	if err != nil {
		return nil, err
	}

	p.platform = platform
	p.application = application

	p.cc = uint32(len(platform.Cores))
	p.tc = uint32(len(application.Tasks))

	p.time = time.NewList(platform, application)
	p.schedule = p.time.Compute(system.NewProfile(platform, application).Mobility)

	C := acorrelation.Compute(application, c.ProbModel.CorrLength)
	p.multiplier, p.zc, err = correlation.Decompose(C, p.tc, c.ProbModel.VarThreshold)
	if err != nil {
		return nil, err
	}

	p.marginals = make([]probability.Inverter, p.tc)
	marginalizer := aprobability.ParseInverter(c.ProbModel.Marginal)
	if marginalizer == nil {
		return nil, errors.New("invalid marginal distributions")
	}
	for i := uint32(0); i < p.tc; i++ {
		duration := platform.Cores[p.schedule.Mapping[i]].Time[application.Tasks[i].Type]
		p.marginals[i] = marginalizer(0, c.ProbModel.MaxDelay*duration)
	}

	return p, nil
}

func (p *Problem) transform(node []float64) []float64 {
	const (
		offset = 1e-8
	)

	z := make([]float64, p.zc)
	u := make([]float64, p.tc)

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
	matrix.Multiply(p.multiplier, z, u, p.tc, p.zc, 1)

	// Dependent Gaussian to dependent uniform to dependent target
	for i := range u {
		u[i] = p.marginals[i].InvCDF(standardGaussian.CDF(u[i]))
	}

	return u
}

func (p *Problem) Printf(format string, arguments ...interface{}) {
	if p.config.Verbose {
		fmt.Printf(format, arguments...)
	}
}

func (p *Problem) Println(arguments ...interface{}) {
	if p.config.Verbose {
		fmt.Println(arguments...)
	}
}
