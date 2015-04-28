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
	Score([]float64, []float64, float64) float64
	Monitor(uint, uint, uint, uint)

	Config() *TargetConfig
}

func NewTarget(problem *Problem) (Target, error) {
	config := problem.Config.Target

	nj, nf, ni := len(config.Rejection), len(config.Refinement), len(config.Importance)
	if nj == 0 || nj != nf || nf != ni {
		return nil, errors.New("the rejection, refinement, and importance " +
			"should not be empty and should have the same number of elements")
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

func Score(target Target, _, surplus []float64, _ float64) float64 {
	config := target.Config()

	rejection, refinement, importance := config.Rejection, config.Refinement, config.Importance

	no, nj := uint(len(surplus)), uint(len(rejection))

	score, reject := 0.0, true
	for i := uint(0); i < no; i++ {
		j := i % nj
		ε := math.Abs(surplus[i])
		if ε >= rejection[j] {
			reject = false
		}
		if ε > refinement[j] {
			score += importance[j] * ε
		}
	}

	if reject {
		score = -1
	}

	return score
}

func Monitor(target Target, k, na, nr, nc uint) {
	if !target.Config().Verbose {
		return
	}
	if k == 0 {
		fmt.Printf("%10s %15s %15s %15s\n", "Step",
			"Accepted Nodes", "Rejected Nodes", "Current Nodes")
	}
	fmt.Printf("%10d %15d %15d %15d\n", k, na, nr, nc)
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
