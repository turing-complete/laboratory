package internal

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

type Target interface {
	Dimensions() (uint, uint)
	Compute([]float64, []float64)
	Refine([]float64, []float64, float64) float64
	Monitor(uint, uint, uint)

	Config() *TargetConfig
}

func NewTarget(problem *Problem) (Target, error) {
	config := problem.Config.Target

	nt, ni := len(config.Tolerance), len(config.Importance)
	if nt == 0 {
		return nil, errors.New("the tolerance should not be empty")
	}
	if ni == 0 {
		return nil, errors.New("the importance should not be empty")
	}
	if nt != ni {
		return nil, errors.New("the tolerance and importance should have the same number of elements")
	}

	switch config.Name {
	case "end-to-end-delay":
		return newDelayTarget(problem, &config), nil
	case "total-energy":
		return newEnergyTarget(problem, &config), nil
	case "temperature-profile":
		return newProfileTarget(problem, &config)
	default:
		return nil, errors.New("the target is unknown")
	}
}

func String(target Target) string {
	ni, no := target.Dimensions()
	return fmt.Sprintf(`{"inputs": %d, "outputs": %d}`, ni, no)
}

func Refine(target Target, _, surplus []float64, _ float64) float64 {
	config := target.Config()

	tolerance, importance := config.Tolerance, config.Importance

	no, nt := uint(len(surplus)), uint(len(tolerance))

	score := 0.0
	for i := uint(0); i < no; i++ {
		j := i % nt
		if w := importance[j]; w > 0 {
			if δ := math.Abs(surplus[i]); δ > tolerance[j] {
				score += w * δ
			}
		}
	}

	return score
}

func Monitor(target Target, k, np, na uint) {
	if !target.Config().Verbose {
		return
	}
	if k == 0 {
		fmt.Printf("%10s %15s %15s\n", "Step", "Passive Nodes", "Active Nodes")
	}
	fmt.Printf("%10d %15d %15d\n", k, np, na)
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
