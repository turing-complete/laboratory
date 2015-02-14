package internal

import (
	"errors"
)

type Target interface {
	Inputs() uint
	Outputs() uint
	Pseudos() uint
	Evaluate([]float64, []float64, []uint64)
	Progress(uint32, uint, uint)
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
