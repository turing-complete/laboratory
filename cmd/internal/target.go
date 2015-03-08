package internal

import (
	"errors"
)

type Target interface {
	Dimensions() (uint, uint)
	Compute([]float64, []float64)
	Refine([]float64) bool
	Monitor(uint, uint, uint)
}

func NewTarget(problem *Problem) (Target, error) {
	config := &problem.Config.Target
	switch config.Name {
	case "end-to-end-delay":
		return newDelayTarget(problem, config), nil
	case "total-energy":
		return newEnergyTarget(problem, config), nil
	case "temperature-profile":
		return newTemperatureTarget(problem, config, &problem.Config.Temperature)
	default:
		return nil, errors.New("the target is unknown")
	}
}
