package internal

import (
	"errors"
	"fmt"

	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"
)

type Target interface {
	Dimensions() (uint, uint)
	Compute([]float64, []float64)
	Refine([]float64) bool
	Monitor(uint, uint, uint)
	Generate(uint) []float64
}

type CommonTarget struct {
	Target
}

func NewTarget(problem *Problem) (Target, error) {
	config := &problem.Config.Target
	switch config.Name {
	case "end-to-end-delay":
		return newDelayTarget(problem, config), nil
	case "total-energy":
		return newEnergyTarget(problem, config), nil
	case "temperature-slice":
		return newSliceTarget(problem, config, &problem.Config.Temperature)
	case "temperature-profile":
		return newProfileTarget(problem, config, &problem.Config.Temperature)
	default:
		return nil, errors.New("the target is unknown")
	}
}

func (t CommonTarget) String() string {
	ni, no := t.Dimensions()
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", ni, no)
}

func (t CommonTarget) Monitor(level, np, na uint) {
	if level == 0 {
		fmt.Printf("%10s %15s %15s\n", "Level", "Passive Nodes", "Active Nodes")
	}
	fmt.Printf("%10d %15d %15d\n", level, np, na)
}

func (t CommonTarget) Generate(ns uint) []float64 {
	ni, _ := t.Dimensions()
	return probability.Sample(uniform.New(0, 1), ns*ni)
}
