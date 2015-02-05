package internal

import (
	"errors"
)

type Target interface {
	Inputs() uint32
	Outputs() uint32
	Evaluate([]float64, []float64, []uint64)
	Progress(uint8, uint32, uint32)
}

func NewTarget(problem *Problem) (Target, error) {
	switch problem.Config.Target {
	case "end-to-end-delay":
		return newDelayTarget(problem), nil
	case "total-energy":
		return newEnergyTarget(problem), nil
	case "temperature-slice":
		return newSliceTarget(problem)
	case "temperature-profile":
		return newProfileTarget(problem)
	default:
		return nil, errors.New("the target is unknown")
	}
}
