package system

import (
	"errors"
	"fmt"
	"log"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/support"
	"github.com/turing-complete/power/dynamic"
	"github.com/turing-complete/power/static"
	"github.com/turing-complete/system"
	"github.com/turing-complete/time"

	temperature "github.com/turing-complete/temperature/analytic"
)

type System struct {
	Platform    *system.Platform
	Application *system.Application

	time         *time.List
	schedule     *time.Schedule
	dynamicPower *dynamic.Power
	staticPower  *static.Power
	temperature  *temperature.Fixed

	Δt float64
}

func New(config *config.System) (*System, error) {
	platform, application, err := system.Load(config.Specification)
	if err != nil {
		return nil, err
	}

	time := time.NewList(platform, application)
	schedule := time.Compute(system.NewProfile(platform, application).Mobility)
	dynamicPower := dynamic.New(platform, application)

	staticPower, err := createStaticPower(dynamicPower, schedule, &config.StaticPower)
	if err != nil {
		return nil, err
	}

	temperature, err := temperature.NewFixed(&config.Config)
	if err != nil {
		return nil, err
	}

	return &System{
		Platform:    platform,
		Application: application,

		time:         time,
		schedule:     schedule,
		dynamicPower: dynamicPower,
		staticPower:  staticPower,
		temperature:  temperature,

		Δt: config.TimeStep,
	}, nil
}

func (self *System) ComputeDynamicPower(schedule *time.Schedule) []float64 {
	return computeDynamicPower(self.dynamicPower, schedule, self.Δt)
}

func (self *System) ComputeSchedule(duration []float64) *time.Schedule {
	return self.time.Update(self.schedule, duration)
}

func (self *System) ComputeTemperatureUpdatePower(P []float64) []float64 {
	nc := uint(self.Platform.Len())
	return self.temperature.ComputeWithStatic(P, func(Q, P []float64) {
		for i := uint(0); i < nc; i++ {
			P[i] += self.staticPower.Compute(Q[i])
		}
	})
}

func (self *System) ReferenceTime() []float64 {
	return self.schedule.Duration()
}

func (self *System) TimeStep() float64 {
	return self.Δt
}

func (self *System) String() string {
	return fmt.Sprintf(`{cores:%d tasks:%d}`,
		self.Platform.Len(), self.Application.Len())
}

func createStaticPower(dynamicPower *dynamic.Power, schedule *time.Schedule,
	config *config.StaticPower) (*static.Power, error) {

	if config.Contribution < 0.0 || config.Contribution >= 1.0 {
		return nil, errors.New("the contribution of the static power is invalid")
	}

	nominal := config.Contribution / (1.0 - config.Contribution) *
		support.Average(dynamicPower.Distribute(schedule))

	if nominal == 0.0 || len(config.Temperature) == 0 || len(config.Coefficient) == 0 {
		log.Printf("The static-power model is disabled.")
		config.Temperature = []float64{0.0, 1.0}
		config.Coefficient = []float64{0.0, 0.0}
	}

	return static.New(nominal, config.Temperature, config.Coefficient), nil
}

func computeDynamicPower(dynamicPower *dynamic.Power,
	schedule *time.Schedule, Δt float64) []float64 {

	return dynamicPower.Sample(schedule, Δt, uint(schedule.Span/Δt))
}
