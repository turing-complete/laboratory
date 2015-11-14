package uncertainty

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

type Uncertainty struct {
	Time  *Parameter
	Power *Parameter
}

func New(system *system.System, config *config.Uncertainty) (*Uncertainty, error) {
	time, err := newParameter(system.ReferenceTime(), config)
	if err != nil {
		return nil, err
	}
	power, err := newParameter(system.ReferencePower(), config)
	if err != nil {
		return nil, err
	}
	return &Uncertainty{
		Time:  time,
		Power: power,
	}, nil
}
