package problem

import (
	"fmt"

	"github.com/simulated-reality/laboratory/src/internal/config"
	"github.com/simulated-reality/laboratory/src/internal/system"
	"github.com/simulated-reality/laboratory/src/internal/uncertainty"
)

type Problem struct {
	Config      *config.Config
	System      *system.System
	Uncertainty *uncertainty.Uncertainty
}

func New(config *config.Config) (*Problem, error) {
	system, err := system.New(&config.System)
	if err != nil {
		return nil, err
	}

	uncertainty, err := uncertainty.New(&config.Uncertainty, system)
	if err != nil {
		return nil, err
	}

	problem := &Problem{
		Config:      config,
		System:      system,
		Uncertainty: uncertainty,
	}

	return problem, nil
}

func (p *Problem) String() string {
	return fmt.Sprintf(`{"system": %s, "uncertainty": %s}`, p.System, p.Uncertainty)
}
