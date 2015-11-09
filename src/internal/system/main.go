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

	return &System{
		Platform:    platform,
		Application: application,

		time:        time,
		power:       power,
		temperature: temperature,

		schedule: schedule,
	}, nil
}

func (self *System) ComputeSchedule(duration []float64) *time.Schedule {
	return self.time.Update(self.schedule, duration)
}

func (self *System) ComputeTime(schedule *time.Schedule) []float64 {
	time := make([]float64, len(schedule.Start))
	for i := range time {
		time[i] = schedule.Finish[i] - schedule.Start[i]
	}
	return time
}

func (self *System) ComputeTemperature(P, ΔT []float64) []float64 {
	return self.temperature.Compute(P, ΔT)
}

func (self *System) DistributePower(schedule *time.Schedule) []float64 {
	cores, tasks := self.Platform.Cores, self.Application.Tasks
	power := make([]float64, self.Application.Len())
	for i, j := range schedule.Mapping {
		power[i] = cores[j].Power[tasks[i].Type]
	}
	return power
}

func (self *System) PartitionPower(schedule *time.Schedule, points []float64,
	ε float64) ([]float64, []float64, []uint) {

	return self.power.Partition(schedule, points, ε)
}

func (self *System) ReferenceTime() []float64 {
	return self.ComputeTime(self.schedule)
}

func (self *System) Span() float64 {
	return self.schedule.Span
}

func (self *System) String() string {
	return fmt.Sprintf(`{"cores": %d, "tasks": %d}`, self.Platform.Len(), self.Application.Len())
}
