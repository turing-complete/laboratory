package main

import (
	"errors"
	"fmt"

	"github.com/ready-steady/persim/system"
	"github.com/ready-steady/persim/time"
	"github.com/ready-steady/prob"
	"github.com/ready-steady/stats/corr"

	"../../pkg/appcorr"
)

type problem struct {
	config Config

	platform    *system.Platform
	application *system.Application

	cc uint32 // cores
	tc uint32 // tasks

	uc uint32 // dependent variables
	zc uint32 // independent variables

	marginals []prob.Inverter
	transform []float64

	time     *time.List
	schedule *time.Schedule
}

func (p *problem) String() string {
	return fmt.Sprintf("Problem{cores: %d, tasks: %d, dvars: %d, ivars: %d}",
		p.cc, p.tc, p.uc, p.zc)
}

func newProblem(config Config) (*problem, error) {
	var err error

	p := &problem{config: config}
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

	C := appcorr.Compute(application, c.TaskIndex, c.ProbModel.CorrLength)
	p.transform, p.zc, err = corr.Decompose(C, p.uc, c.ProbModel.VarThreshold)
	if err != nil {
		return nil, err
	}

	p.marginals = make([]prob.Inverter, p.uc)
	marginalizer := marginalize(c.ProbModel.Marginal)
	if marginalizer == nil {
		return nil, errors.New("invalid marginal distributions")
	}
	for i, tid := range c.TaskIndex {
		duration := platform.Cores[p.schedule.Mapping[tid]].Time[application.Tasks[tid].Type]
		p.marginals[i] = marginalizer(c.ProbModel.MaxDelay * duration)
	}

	return p, nil
}
