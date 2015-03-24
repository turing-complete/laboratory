package internal

import (
	"errors"
	"fmt"
	"math"

	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"
)

type Target interface {
	Config() *TargetConfig
	Dimensions() (uint, uint)

	Compute([]float64, []float64)
	Refine([]float64, []float64, float64) float64
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
	case "temperature-switch":
		return newSwitchTarget(problem, config, &problem.Config.Temperature)
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

func (t CommonTarget) Refine(_, surplus []float64, volume float64) float64 {
	nm := uint(len(surplus)) / 2

	k := uint(0)
	if t.Config().Squared {
		k = 1
	}

	Σ := 0.0
	for i := uint(0); i < nm; i++ {
		Δ := surplus[i*2+k] * volume
		Σ += Δ * Δ
	}
	Σ = math.Sqrt(Σ)

	if Σ <= t.Config().Tolerance {
		Σ = 0
	}

	return Σ
}

func (t CommonTarget) Monitor(k, np, na uint) {
	if !t.Config().Verbose {
		return
	}
	if k == 0 {
		fmt.Printf("%10s %15s %15s\n", "Iteration", "Passive Nodes", "Active Nodes")
	}
	fmt.Printf("%10d %15d %15d\n", k, np, na)
}

func (t CommonTarget) Generate(ns uint) []float64 {
	ni, _ := t.Dimensions()
	return probability.Sample(uniform.New(0, 1), ns*ni)
}
