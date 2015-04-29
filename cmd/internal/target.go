package internal

import (
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/ready-steady/adapt"
)

type Target adapt.Target

func NewTarget(problem *Problem) (Target, error) {
	config := problem.Config.Target

	nj, nf := len(config.Rejection), len(config.Refinement)
	if nj == 0 || nj != nf {
		return nil, errors.New("the rejection and refinement " +
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

func Monitor(target Target, progress *adapt.Progress) {
	if progress.Iteration == 0 {
		fmt.Printf("%10s %15s %15s %15s\n", "Iteration",
			"Accepted Nodes", "Rejected Nodes", "Current Nodes")
	}
	fmt.Printf("%10d %15d %15d %15d\n", progress.Iteration,
		progress.Accepted, progress.Rejected, progress.Current)
}

func Score(target Target, config *TargetConfig,
	location *adapt.Location, progress *adapt.Progress) float64 {

	rejection := config.Rejection
	refinement := config.Refinement

	_, no := target.Dimensions()
	nj := uint(len(rejection))

	score, reject := 0.0, true
	for i := uint(0); i < no; i++ {
		j := i % nj
		ε := math.Abs(location.Surplus[i])
		if ε >= rejection[j] {
			reject = false
		}
		if ε > refinement[j] {
			score += ε
		}
	}

	if reject {
		score = -1
	}

	return score
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
