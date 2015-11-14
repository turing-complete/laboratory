package uncertainty

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

type Parameter interface {
	Dimensions() uint
	Transform([]float64) []float64
}

type Uncertainty struct {
	Time  Parameter
	Power Parameter
}

func New(system *system.System, config *config.Uncertainty) (*Uncertainty, error) {
	time, err := newDirect(system.ReferenceTime(), config)
	if err != nil {
		return nil, err
	}
	power, err := newDirect(system.ReferencePower(), config)
	if err != nil {
		return nil, err
	}
	return &Uncertainty{
		Time:  time,
		Power: power,
	}, nil
}
