package internal

import (
	"fmt"

	"github.com/ready-steady/simulation/power"
	asystem "github.com/ready-steady/simulation/system"
	"github.com/ready-steady/simulation/temperature/analytic"
	"github.com/ready-steady/simulation/time"
)

type system struct {
	platform    *asystem.Platform
	application *asystem.Application

	time        *time.List
	power       *power.Power
	temperature *analytic.Fluid

	schedule  *time.Schedule
	reference []float64

	nc uint
	nt uint
}

func newSystem(config *SystemConfig) (*system, error) {
	platform, application, err := asystem.Load(config.Specification)
	if err != nil {
		return nil, err
	}

	time := time.NewList(platform, application)
	power := power.New(platform, application)
	temperature, err := analytic.NewFluid(&config.Config)
	if err != nil {
		return nil, err
	}

	schedule := time.Compute(asystem.NewProfile(platform, application).Mobility)

	nc := uint(len(platform.Cores))
	nt := uint(len(application.Tasks))

	system := &system{
		platform:    platform,
		application: application,

		time:        time,
		power:       power,
		temperature: temperature,

		schedule:  schedule,
		reference: computeTime(schedule),

		nc: nc,
		nt: nt,
	}

	return system, nil
}

func (s *system) String() string {
	return fmt.Sprintf(`{"cores": %d, "tasks": %d}`, s.nc, s.nt)
}

func (s *system) computeSchedule(modes []float64) *time.Schedule {
	duration := make([]float64, s.nt)
	for i, time := range s.reference {
		duration[i] = (1 + modes[i]) * time
	}
	return s.time.Update(s.schedule, duration)
}

func (s *system) computeTime(schedule *time.Schedule) []float64 {
	return computeTime(schedule)
}

func (s *system) computePower(schedule *time.Schedule) []float64 {
	cores, tasks := s.platform.Cores, s.application.Tasks

	power := make([]float64, s.nt)
	for i, j := range schedule.Mapping {
		power[i] = cores[j].Power[tasks[i].Type]
	}

	return power
}

func computeTime(schedule *time.Schedule) []float64 {
	time := make([]float64, len(schedule.Start))
	for i := range time {
		time[i] = schedule.Finish[i] - schedule.Start[i]
	}
	return time
}
