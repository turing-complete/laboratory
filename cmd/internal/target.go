package internal

import (
	"errors"
)

type Target interface {
	Inputs() uint
	Outputs() uint
	Evaluate([]float64, []float64, []uint64)
	Progress(uint32, uint, uint)
}

func NewTarget(problem *Problem) (Target, error) {
	switch problem.Config.Target.Name {
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
