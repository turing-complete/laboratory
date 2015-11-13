package uncertainty

import (
	"github.com/turing-complete/laboratory/src/internal/config"
)

type Uncertainty interface {
	Parameters() uint
	Transform([]float64) []float64
}

func New(reference []float64, config *config.Uncertainty) (Uncertainty, error) {
	return newDirect(reference, config)
}
