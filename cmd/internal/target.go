package internal

import (
	"errors"

	"../../pkg/solver"
)

type Target interface {
	InputsOutputs() (uint32, uint32)
	Serve(<-chan solver.Job)
}

func newTarget(p *Problem) (Target, error) {
	switch p.config.Target {
	case "end-to-end-delay":
		return newSpanTarget(p)
	case "temperature-profile":
		return newTempTarget(p)
	default:
		return nil, errors.New("the target is unknown")
	}
}
