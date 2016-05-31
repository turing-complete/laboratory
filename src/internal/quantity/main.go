package quantity

import (
	"errors"

	"github.com/ready-steady/lapack"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"

	interpolation "github.com/ready-steady/adapt/algorithm"
)

func init() {
	// The quantities of interest involve linear algebra, which is powered by
	// OpenBLAS via the lapack package. They are evaluated in multiple threads;
	// however, OpenBLAS is multithreaded by itself. The two multithreading
	// implementations might collide. Hence, the OpenBLAS one must be disabled.
	lapack.SetNumberOfThreads(1)
}

type Quantity interface {
	Dimensions() (uint, uint)
	Compute([]float64, []float64)

	Evaluate([]float64) float64
	Forward([]float64) []float64
	Backward([]float64) []float64
}

func New(system *system.System, uncertainty uncertainty.Uncertainty,
	config *config.Quantity) (Quantity, error) {

	switch config.Name {
	case "end-to-end-delay":
		return newDelay(system, uncertainty, config)
	case "total-energy":
		return newEnergy(system, uncertainty, config)
	case "maximum-temperature":
		return newTemperature(system, uncertainty, config)
	default:
		return nil, errors.New("the quantity is unknown")
	}
}

func Invoke(quantity Quantity, points []float64) []float64 {
	ni, no := quantity.Dimensions()
	return interpolation.Invoke(quantity.Compute, points, ni, no)
}
