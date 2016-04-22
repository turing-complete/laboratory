package uncertainty

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

func newEpistemic(system *system.System, reference []float64,
	config *config.Parameter) (*aleatory, error) {

	config.Distribution = "Uniform()"
	config.Correlation = 0.0
	return newAleatory(system, reference, config)
}
