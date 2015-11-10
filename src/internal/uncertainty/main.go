package uncertainty

import (
	"errors"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

type Uncertainty interface {
	Len() int
	Transform([]float64) []float64
}

func New(system *system.System, reference []float64,
	config *config.Uncertainty) (Uncertainty, error) {

	switch config.Name {
	case "direct":
		return newDirect(system, reference, config)
	case "marginal":
		return newMarginal(system, reference, config)
	default:
		return nil, errors.New("the uncertainty model is unknown")
	}
}
