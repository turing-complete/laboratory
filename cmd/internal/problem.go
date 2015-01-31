package internal

import (
	"errors"
	"fmt"

	"github.com/ready-steady/probability"
	"github.com/ready-steady/simulation/system"
	"github.com/ready-steady/simulation/time"
	"github.com/ready-steady/statistics/correlation"

	acorrelation "../../pkg/correlation"
	aprobability "../../pkg/probability"
	"../../pkg/solver"
)

type Problem struct {
	config Config

	platform    *system.Platform
	application *system.Application

	cc uint32 // cores
	tc uint32 // tasks

	uc uint32 // dependent variables
	zc uint32 // independent variables

	marginals []probability.Inverter
	transform []float64

	time     *time.List
	schedule *time.Schedule
}

func (p *Problem) String() string {
	return fmt.Sprintf("Problem{cores: %d, tasks: %d, dvars: %d, ivars: %d}",
		p.cc, p.tc, p.uc, p.zc)
}

func newProblem(config Config) (*Problem, error) {
	var err error

	p := &Problem{config: config}
	c := &p.config

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

	p.uc = p.tc

	C := acorrelation.Compute(application, c.ProbModel.CorrLength)
	p.transform, p.zc, err = correlation.Decompose(C, p.uc, c.ProbModel.VarThreshold)
	if err != nil {
		return nil, err
	}

	p.marginals = make([]probability.Inverter, p.uc)
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

func (p *Problem) Setup() (Target, *solver.Solver, error) {
	target, err := newTarget(p)
	if err != nil {
		return nil, nil, err
	}

	ic, oc := target.InputsOutputs()

	config := p.config.Solver
	config.Inputs = uint16(ic)
	config.Outputs = uint16(oc)
	config.ArtificialInputs = uint16(ic - p.zc)

	solver, err := solver.New(config, target.Serve)
	if err != nil {
		return nil, nil, err
	}

	return target, solver, nil
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
