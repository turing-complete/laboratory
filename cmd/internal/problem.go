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

	p.time = time.NewList(platform, application)
	p.schedule = p.time.Compute(system.NewProfile(platform, application).Mobility)

	p.uc = uint32(len(c.TaskIndex))

	C := acorrelation.Compute(application, c.TaskIndex, c.ProbModel.CorrLength)
	p.transform, p.zc, err = correlation.Decompose(C, p.uc, c.ProbModel.VarThreshold)
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

func (p *Problem) Setup() (Target, *solver.Solver, error) {
	target, err := newTarget(p)
	if err != nil {
		return nil, nil, err
	}

	ic, oc := target.InputsOutputs()

	config := p.config.Solver
	config.Inputs = uint16(ic)
	config.Outputs = uint16(oc)
	config.CacheInputs = uint16(ic - p.zc)

	solver, err := solver.New(config, target.Serve)
	if err != nil {
		return nil, nil, err
	}

	return target, solver, nil
}

func (p *Problem) Log(a ...interface{}) {
	if p.config.Verbose {
		fmt.Println(a...)
	}
}
