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
		return newDelayTarget(p)
	case "total-energy":
		return newEnergyTarget(p)
	case "temperature-profile":
		return newHeatTarget(p)
	default:
		return nil, errors.New("the target is unknown")
	}
}

func processNode(node float64) float64 {
	const (
		offset = 1e-8
	)

	switch node {
	case 0:
		node = 0 + offset
	case 1:
		node = 1 - offset
	}

	return node
}
