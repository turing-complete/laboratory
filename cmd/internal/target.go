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

	if len(config.Importance) == 0 {
		return nil, errors.New("the importance should not be empty")
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

	_, no := target.Dimensions()

	importance := config.Importance
	ni := uint(len(importance))

	score, norm := 0.0, 0.0
	for i := uint(0); i < no; i++ {
		α := importance[i%ni]
		if α == 0 {
			continue
		}

		s := α * location.Volume * location.Surplus[i]
		score += s * s

		n := α * progress.Integral[i]
		norm += n * n
	}

	score, norm = math.Sqrt(score), math.Sqrt(norm)
	if norm > 0 {
		score /= norm
	}

	if score < config.Rejection {
		return -1
	}
	if score < config.Refinement {
		return 0
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
