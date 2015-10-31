package system

import (
	"fmt"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/power"
	"github.com/turing-complete/system"
	"github.com/turing-complete/temperature/analytic"
	"github.com/turing-complete/time"
)

type System struct {
	Platform    *system.Platform
	Application *system.Application

	time        *time.List
	power       *power.Power
	temperature *analytic.Fluid

	schedule *time.Schedule
}

func New(config *config.System) (*System, error) {
	platform, application, err := system.Load(config.Specification)
	if err != nil {
		return nil, err
	}

	time := time.NewList(platform, application)
	power := power.New(platform, application)
	temperature, err := analytic.NewFluid(&config.Config)
	if err != nil {
		return nil, err
	}

	schedule := time.Compute(system.NewProfile(platform, application).Mobility)

	system := &System{
		Platform:    platform,
		Application: application,

		time:        time,
		power:       power,
		temperature: temperature,

		schedule: schedule,
	}

	return system, nil
}

func (s *System) ComputeSchedule(duration []float64) *time.Schedule {
	return s.time.Update(s.schedule, duration)
}

func (s *System) ComputeTime(schedule *time.Schedule) []float64 {
	return computeTime(schedule)
}

func (s *System) ComputeTemperature(P, ΔT []float64) []float64 {
	return s.temperature.Compute(P, ΔT)
}

func (s *System) DistributePower(schedule *time.Schedule) []float64 {
	cores, tasks := s.Platform.Cores, s.Application.Tasks

	power := make([]float64, s.Application.Len())
	for i, j := range schedule.Mapping {
		power[i] = cores[j].Power[tasks[i].Type]
	}

	return power
}

func (s *System) PartitionPower(schedule *time.Schedule, points []float64,
	ε float64) ([]float64, []float64, []uint) {

	return s.power.Partition(schedule, points, ε)
}

func (s *System) ReferenceTime() []float64 {
	return computeTime(s.schedule)
}

func (s *System) Span() float64 {
	return s.schedule.Span
}

func (s *System) String() string {
	return fmt.Sprintf(`{"cores": %d, "tasks": %d}`, s.Platform.Len(), s.Application.Len())
}

func computeTime(schedule *time.Schedule) []float64 {
	time := make([]float64, len(schedule.Start))
	for i := range time {
		time[i] = schedule.Finish[i] - schedule.Start[i]
	}
	return time
}
