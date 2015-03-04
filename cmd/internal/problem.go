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

	time     *time.List
	schedule *time.Schedule

	tasks      []uint
	multiplier []float64
	marginals  []probability.Inverter

	nc uint
	nt uint
	nu uint
	nz uint
}

func (p *Problem) String() string {
	return fmt.Sprintf("Problem{cores: %d, tasks: %d, dvars: %d, ivars: %d}",
		p.nc, p.nt, p.nu, p.nz)
}

func NewProblem(config Config) (*Problem, error) {
	platform, application, err := system.Load(config.TGFF)
	if err != nil {
		return nil, err
	}

	nc, nt := uint(len(platform.Cores)), uint(len(application.Tasks))

	time := time.NewList(platform, application)
	schedule := time.Compute(system.NewProfile(platform, application).Mobility)

	c := &config.Probability

	if c.MaxDelay < 0 || 1 <= c.MaxDelay {
		return nil, errors.New("the delay rate is invalid")
	}
	if c.CorrLength <= 0 {
		return nil, errors.New("the correlation length is invalid")
	}
	if c.VarThreshold <= 0 {
		return nil, errors.New("the variance-reduction threshold is invalid")
	}

	tasks := c.TaskIndex
	if len(tasks) == 0 {
		tasks = make([]uint, nt)
		for i := uint(0); i < nt; i++ {
			tasks[i] = i
		}
	}

	nu := uint(len(tasks))

	C := acorrelation.Compute(application, tasks, c.CorrLength)
	multiplier, nz, err := correlation.Decompose(C, nu, c.VarThreshold)
	if err != nil {
		return nil, err
	}

	marginalizer := aprobability.ParseInverter(c.Marginal)
	if marginalizer == nil {
		return nil, errors.New("invalid marginal distributions")
	}
	marginals := make([]probability.Inverter, nu)
	for i, tid := range tasks {
		duration := platform.Cores[schedule.Mapping[tid]].Time[application.Tasks[tid].Type]
		marginals[i] = marginalizer(0, c.MaxDelay*duration)
	}

	problem := &Problem{
		Config: config,

		platform:    platform,
		application: application,

		time:     time,
		schedule: schedule,

		tasks:      tasks,
		multiplier: multiplier,
		marginals:  marginals,

		nc: nc,
		nt: nt,
		nu: nu,
		nz: nz,
	}

	return problem, nil
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
	for i, tid := range p.tasks {
		d[tid] = p.marginals[i].InvCDF(standardGaussian.CDF(u[i]))
	}

	return d
}
