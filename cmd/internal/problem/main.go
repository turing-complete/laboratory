package problem

import (
	"fmt"

	"github.com/simulated-reality/laboratory/cmd/internal/config"
	"github.com/simulated-reality/laboratory/cmd/internal/model"
	"github.com/simulated-reality/laboratory/cmd/internal/system"
)

type Problem struct {
	Config *config.Config
	System *system.System
	Model  *model.Model
}

func (p *Problem) String() string {
	return fmt.Sprintf(`{"system": %s, "model": %s}`, p.System, p.Model)
}

func New(config *config.Config) (*Problem, error) {
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
		System: system,
		Model:  model,
	}

	return problem, nil
}
