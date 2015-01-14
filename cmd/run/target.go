package main

import (
	"errors"

	"../../pkg/solver"
)

type target interface {
	InputsOutputs() (uint32, uint32)
	Serve(<-chan solver.Job)
}

func newTarget(p *problem) (target, error) {
	switch p.config.Target {
	case "end-to-end-delay":
		return newSpanTarget(p)
	case "temperature-profile":
		return newTempTarget(p)
	default:
		return nil, errors.New("the target is unknown")
	}
}
