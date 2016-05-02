package uncertainty

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

type Uncertainty interface {
	Mapping() (uint, uint)
	Forward([]float64) []float64
	Backward([]float64) []float64
	Weigh([]float64) float64
}

func NewAleatory(system *system.System, config *config.Uncertainty) (Uncertainty, error) {
	return newBase(system, system.ReferenceTime(), config)
}

func NewEpistemic(system *system.System, config *config.Uncertainty) (Uncertainty, error) {
	config.Distribution, config.Correlation, config.Variance = "Uniform()", 0.0, 1.0
	return newBase(system, system.ReferenceTime(), config)
}
