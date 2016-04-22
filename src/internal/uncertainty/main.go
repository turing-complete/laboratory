package uncertainty

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

type Transform interface {
	Mapping() (uint, uint)
	Forward([]float64) []float64
	Inverse([]float64) []float64
}

type Uncertainty struct {
	Time  Transform
	Power Transform
}

func NewAleatory(system *system.System, config *config.Uncertainty) (*Uncertainty, error) {
	time, err := newAleatory(system, system.ReferenceTime(), &config.Time)
	if err != nil {
		return nil, err
	}
	power, err := newAleatory(system, system.ReferencePower(), &config.Power)
	if err != nil {
		return nil, err
	}
	return &Uncertainty{
		Time:  time,
		Power: power,
	}, nil
}

func NewEpistemic(system *system.System, config *config.Uncertainty) (*Uncertainty, error) {
	time, err := newEpistemic(system, system.ReferenceTime(), &config.Time)
	if err != nil {
		return nil, err
	}
	power, err := newEpistemic(system, system.ReferencePower(), &config.Power)
	if err != nil {
		return nil, err
	}
	return &Uncertainty{
		Time:  time,
		Power: power,
	}, nil
}

func (self *Uncertainty) Mapping() (uint, uint) {
	ni1, no1 := self.Time.Mapping()
	ni2, no2 := self.Power.Mapping()
	return ni1 + ni2, no1 + no2
}

func (self *Uncertainty) Forward(ω []float64) []float64 {
	_, no := self.Time.Mapping()
	return append(self.Time.Forward(ω[:no]), self.Power.Forward(ω[no:])...)
}

func (self *Uncertainty) Inverse(z []float64) []float64 {
	ni, _ := self.Time.Mapping()
	return append(self.Time.Inverse(z[:ni]), self.Power.Inverse(z[ni:])...)
}
