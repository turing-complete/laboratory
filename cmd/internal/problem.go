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
	uc uint32
	zc uint32

	marginals  []probability.Inverter
	multiplier []float64

	time     *time.List
	schedule *time.Schedule
}

func (p *Problem) String() string {
	return fmt.Sprintf("Problem{cores: %d, tasks: %d, dvars: %d, ivars: %d}",
		p.cc, p.tc, p.uc, p.zc)
}

func NewProblem(config string) (*Problem, error) {
	c, err := loadConfig(config)
	if err != nil {
		return nil, err
	}

	if c.ProbModel.MaxDelay < 0 || 1 <= c.ProbModel.MaxDelay {
		return nil, errors.New("the delay rate is invalid")
	}
	if c.ProbModel.CorrLength <= 0 {
		return nil, errors.New("the correlation length is invalid")
	}
	if c.ProbModel.VarThreshold <= 0 || 1 < c.ProbModel.VarThreshold {
		return nil, errors.New("the variance-reduction threshold is invalid")
	}

	p := &Problem{config: c}

	platform, application, err := system.Load(c.TGFF)
	if err != nil {
		return nil, err
	}

	p.platform = platform
	p.application = application

	p.cc = uint32(len(platform.Cores))
	p.tc = uint32(len(application.Tasks))

	if len(c.CoreIndex) == 0 {
		c.CoreIndex = make([]uint16, p.cc)
		for i := uint16(0); i < uint16(p.cc); i++ {
			c.CoreIndex[i] = i
		}
	}
	if len(c.TaskIndex) == 0 {
		c.TaskIndex = make([]uint16, p.tc)
		for i := uint16(0); i < uint16(p.tc); i++ {
			c.TaskIndex[i] = i
		}
	}

	p.uc = uint32(len(c.TaskIndex))

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
	for i, tid := range p.config.TaskIndex {
		d[tid] = p.marginals[i].InvCDF(standardGaussian.CDF(u[i]))
	}

	return d
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
