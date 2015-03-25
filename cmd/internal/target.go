package internal

import (
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"
)

type Target interface {
	Dimensions() (uint, uint)
	Compute([]float64, []float64)
	Refine([]float64, []float64, float64) float64
	Monitor(uint, uint, uint)

	Config() *TargetConfig
	Generate(uint) []float64
}

type GenericTarget struct {
	Target
}

func NewTarget(problem *Problem) (Target, error) {
	config := problem.Config.Target

	if len(config.Stencil) == 0 {
		config.Stencil = []bool{true, false}
	}

	switch config.Name {
	case "end-to-end-delay":
		return newDelayTarget(problem, &config), nil
	case "total-energy":
		return newEnergyTarget(problem, &config), nil
	case "temperature-slice":
		return newSliceTarget(problem, &config, &problem.Config.Temperature)
	case "temperature-switch":
		return newSwitchTarget(problem, &config, &problem.Config.Temperature)
	case "temperature-profile":
		return newProfileTarget(problem, &config, &problem.Config.Temperature)
	default:
		return nil, errors.New("the target is unknown")
	}
}

func (t GenericTarget) String() string {
	ni, no := t.Dimensions()
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", ni, no)
}

func (t GenericTarget) Refine(_, surplus []float64, volume float64) float64 {
	config := t.Config()

	stencil := config.Stencil

	no, ns := uint(len(surplus)), uint(len(stencil))

	Σ := 0.0
	for i := uint(0); i < no; i++ {
		if stencil[i%ns] {
			s := surplus[i] * volume
			Σ += s * s
		}
	}
	Σ = math.Sqrt(Σ)

	if Σ <= config.Tolerance {
		Σ = 0
	}

	return Σ
}

func (t GenericTarget) Monitor(k, np, na uint) {
	if !t.Config().Verbose {
		return
	}
	if k == 0 {
		fmt.Printf("%10s %15s %15s\n", "Iteration", "Passive Nodes", "Active Nodes")
	}
	fmt.Printf("%10d %15d %15d\n", k, np, na)
}

func (t GenericTarget) Generate(ns uint) []float64 {
	ni, _ := t.Dimensions()
	return probability.Sample(uniform.New(0, 1), ns*ni)
}

func Invoke(target Target, points []float64, nw uint) []float64 {
	ni, no := target.Dimensions()
	np := uint(len(points)) / ni

	values := make([]float64, np*no)
	jobs := make(chan uint, np)
	group := sync.WaitGroup{}
	group.Add(int(np))

	for i := uint(0); i < nw; i++ {
		go func() {
			for j := range jobs {
				target.Compute(points[j*ni:(j+1)*ni], values[j*no:(j+1)*no])
				group.Done()
			}
		}()
	}

	for i := uint(0); i < np; i++ {
		jobs <- i
	}

	group.Wait()
	close(jobs)

	return values
}
