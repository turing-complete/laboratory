package internal

import (
	"fmt"
)

type Problem struct {
	Config *Config
	system *system
	model  *model
}

func (p *Problem) String() string {
	return fmt.Sprintf(`{"system": %s, "model": %s}`, p.system, p.model)
}

func NewProblem(config *Config) (*Problem, error) {
	system, err := newSystem(&config.System)
	if err != nil {
		return nil, err
	}

	model, err := newModel(&config.Probability, system)
	if err != nil {
		return nil, err
	}

	problem := &Problem{
		Config: config,
		system: system,
		model:  model,
	}

	return problem, nil
}
