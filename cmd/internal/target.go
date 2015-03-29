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

func String(target Target) string {
	ni, no := target.Dimensions()
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", ni, no)
}

func Refine(target Target, _, surplus []float64, volume float64) float64 {
	config := target.Config()

	stencil := config.Stencil

	no, ns := uint(len(surplus)), uint(len(stencil))

	score := 0.0
	for i := uint(0); i < no; i++ {
		if stencil[i%ns] {
			score += surplus[i] * surplus[i]
		}
	}
	score = volume * math.Sqrt(score)

	if score <= config.Tolerance {
		score = 0
	}

	return score
}

func Monitor(target Target, k, np, na uint) {
	if !target.Config().Verbose {
		return
	}
	if k == 0 {
		fmt.Printf("%10s %15s %15s\n", "Iteration", "Passive Nodes", "Active Nodes")
	}
	fmt.Printf("%10d %15d %15d\n", k, np, na)
}

func Generate(target Target, ns uint) []float64 {
	ni, _ := target.Dimensions()
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
