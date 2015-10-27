package internal

import (
	"fmt"

	"github.com/simulated-reality/laboratory/internal/config"
	"github.com/simulated-reality/laboratory/internal/model"
	"github.com/simulated-reality/laboratory/internal/system"
)

type Problem struct {
	Config *config.Config
	system *system.System
	model  *model.Model
}

func (p *Problem) String() string {
	return fmt.Sprintf(`{"system": %s, "model": %s}`, p.system, p.model)
}

func NewProblem(config *config.Config) (*Problem, error) {
	system, err := system.New(&config.System)
	if err != nil {
		return nil, err
	}

	model, err := model.New(&config.Probability, system)
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
