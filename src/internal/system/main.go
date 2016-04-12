package system

import (
	"fmt"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/power"
	"github.com/turing-complete/system"
	"github.com/turing-complete/time"

	temperature "github.com/turing-complete/temperature/analytic"
)

type System struct {
	Platform    *system.Platform
	Application *system.Application

	time        *time.List
	power       *power.Power
	temperature *temperature.Fluid

	schedule *time.Schedule
}

func New(config *config.System) (*System, error) {
	platform, application, err := system.Load(config.Specification)
	if err != nil {
		return nil, err
	}

	time := time.NewList(platform, application)
	power := power.New(platform, application)
	temperature, err := temperature.NewFluid(&config.Config)
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

func (self *System) ComputeTemperature(P, ΔT []float64) []float64 {
	return self.temperature.Compute(P, ΔT)
}

func (self *System) PartitionPower(values []float64, schedule *time.Schedule,
	ε float64) ([]float64, []float64) {

	return power.Partition(values, schedule, ε)
}

func (self *System) ReferencePower() []float64 {
	return self.power.Distribute(self.schedule)
}

func (self *System) ReferenceTime() []float64 {
	return self.schedule.Duration()
}

func (self *System) Span() float64 {
	return self.schedule.Span
}

func (self *System) String() string {
	return fmt.Sprintf(`{cores:%d tasks:%d}`, self.Platform.Len(), self.Application.Len())
}
