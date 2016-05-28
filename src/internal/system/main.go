package system

import (
	"fmt"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/power/dynamic"
	"github.com/turing-complete/system"
	"github.com/turing-complete/time"

	temperature "github.com/turing-complete/temperature/analytic"
)

type System struct {
	Platform    *system.Platform
	Application *system.Application

	time        *time.List
	schedule    *time.Schedule
	power       *dynamic.Power
	temperature *temperature.Fixed

	Δt float64
}

func New(config *config.System) (*System, error) {
	platform, application, err := system.Load(config.Specification)
	if err != nil {
		return nil, err
	}

	time := time.NewList(platform, application)
	schedule := time.Compute(system.NewProfile(platform, application).Mobility)
	power := dynamic.New(platform, application)
	temperature, err := temperature.NewFixed(&config.Config)
	if err != nil {
		return nil, err
	}

	return &System{
		Platform:    platform,
		Application: application,

		time:        time,
		schedule:    schedule,
		power:       power,
		temperature: temperature,

		Δt: config.TimeStep,
	}, nil
}

func (self *System) ComputePower(schedule *time.Schedule) []float64 {
	return self.power.Sample(schedule, self.Δt, uint(schedule.Span/self.Δt))
}

func (self *System) ComputeSchedule(duration []float64) *time.Schedule {
	return self.time.Update(self.schedule, duration)
}

func (self *System) ComputeTemperature(P []float64) []float64 {
	return self.temperature.Compute(P)
}

func (self *System) ReferenceTime() []float64 {
	return self.schedule.Duration()
}

func (self *System) TimeStep() float64 {
	return self.Δt
}

func (self *System) String() string {
	return fmt.Sprintf(`{cores:%d tasks:%d}`, self.Platform.Len(), self.Application.Len())
}
