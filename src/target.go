package main

import (
	"errors"
)

type job struct {
	key   string
	data  []float64
	node  []float64
	value []float64
	done  chan<- result
}

type result struct {
	key  string
	data []float64
}

type target interface {
	InputsOutputs() (uint32, uint32)
	Serve(<-chan job)
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
