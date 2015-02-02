package internal

import (
	"errors"
)

type Target interface {
	Evaluate([]float64, []float64, []uint64)
	InputsOutputs() (uint32, uint32)
	Evaluations() uint32
}

func NewTarget(problem *Problem) (Target, error) {
	switch problem.config.Target {
	case "end-to-end-delay":
		return newDelayTarget(problem), nil
	case "total-energy":
		return newEnergyTarget(problem), nil
	case "temperature-profile":
		return newTemperatureTarget(problem)
	default:
		return nil, errors.New("the target is unknown")
	}
}
