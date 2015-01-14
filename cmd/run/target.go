package main

import (
	"errors"
)

type target interface {
	InputsOutputs() (uint32, uint32)
	Serve(<-chan job)
}

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
