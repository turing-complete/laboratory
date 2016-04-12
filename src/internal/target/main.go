package target

import (
	"errors"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"

	interpolation "github.com/ready-steady/adapt/algorithm/external"
)

type Target interface {
	Dimensions() (uint, uint)
	Compute([]float64, []float64)
	Forward([]float64) []float64
	Inverse([]float64) []float64
}

func New(system *system.System, uncertainty *uncertainty.Uncertainty,
	config *config.Target) (Target, error) {

	switch config.Name {
	case "end-to-end-delay":
		return newDelay(system, uncertainty, config)
	case "total-energy":
		return newEnergy(system, uncertainty, config)
	case "maximal-temperature":
		return newTemperature(system, uncertainty, config)
	default:
		return nil, errors.New("the target is unknown")
	}
}

func Invoke(target Target, points []float64) []float64 {
	ni, no := target.Dimensions()
	return interpolation.Invoke(target.Compute, points, ni, no)
}
