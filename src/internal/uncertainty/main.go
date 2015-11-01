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

func New(system *system.System, config *config.Uncertainty) (Uncertainty, error) {
	switch config.Name {
	case "direct":
		return newDirect(system, config)
	case "marginal":
		return newMarginal(system, config)
	default:
		return nil, errors.New("the uncertainty model is unknown")
	}
}
