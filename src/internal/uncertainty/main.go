package uncertainty

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

type Uncertainty interface {
	Len() int
	Transform([]float64) []float64
}

func New(s *system.System, c *config.Uncertainty) (Uncertainty, error) {
	return newMarginal(s, c)
}
