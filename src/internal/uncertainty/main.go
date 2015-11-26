package uncertainty

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

type Parameter interface {
	Dimensions() (uint, uint)
	Forward([]float64) []float64
	Inverse([]float64) []float64
}

type Uncertainty struct {
	Time  Parameter
	Power Parameter
}

func New(system *system.System, config *config.Uncertainty) (*Uncertainty, error) {
	time, err := newDirect(system.ReferenceTime(), &config.Time)
	if err != nil {
		return nil, err
	}
	power, err := newDirect(system.ReferencePower(), &config.Power)
	if err != nil {
		return nil, err
	}
	return &Uncertainty{
		Time:  time,
		Power: power,
	}, nil
}

func NewMarginal(system *system.System, config *config.Uncertainty) (*Uncertainty, error) {
	time, err := newMarginal(system, system.ReferenceTime(), &config.Time)
	if err != nil {
		return nil, err
	}
	power, err := newMarginal(system, system.ReferencePower(), &config.Power)
	if err != nil {
		return nil, err
	}
	return &Uncertainty{
		Time:  time,
		Power: power,
	}, nil
}

func (self *Uncertainty) Dimensions() (uint, uint) {
	ni1, no1 := self.Time.Dimensions()
	ni2, no2 := self.Power.Dimensions()
	return ni1 + ni2, no1 + no2
}

func (self *Uncertainty) Forward(z []float64) []float64 {
	ni, _ := self.Time.Dimensions()
	return append(self.Time.Forward(z[:ni]), self.Power.Forward(z[ni:])...)
}

func (self *Uncertainty) Inverse(ω []float64) []float64 {
	_, no := self.Time.Dimensions()
	return append(self.Time.Inverse(ω[:no]), self.Power.Inverse(ω[no:])...)
}
