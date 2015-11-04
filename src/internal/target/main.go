package target

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/ready-steady/adapt"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

type Target adapt.Target

func New(system *system.System, config *config.Target) (Target, error) {
	ni, nj, nf := len(config.Importance), len(config.Rejection), len(config.Refinement)
	if ni == 0 || ni != nj || nj != nf {
		return nil, errors.New("the importance, refinement, and rejection " +
			"should not be empty and should have the same number of elements")
	}

	switch config.Name {
	case "end-to-end-delay":
		return newDelay(system, config)
	case "total-energy":
		return newEnergy(system, config)
	case "temperature-profile":
		return newProfile(system, config)
	default:
		return nil, errors.New("the target is unknown")
	}
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

func display(target Target) string {
	ni, no := target.Dimensions()
	return fmt.Sprintf(`{"inputs": %d, "outputs": %d}`, ni, no)
}

func monitor(target Target, progress *adapt.Progress) {
	if progress.Iteration == 0 {
		log.Printf("%5s %10s %15s %15s %15s\n", "Level", "Iteration",
			"Accepted Nodes", "Rejected Nodes", "Current Nodes")
	}
	log.Printf("%5d %10d %15d %15d %15d\n", progress.Level, progress.Iteration,
		progress.Accepted, progress.Rejected, progress.Current)
}

func score(target Target, config *config.Target, location *adapt.Location,
	progress *adapt.Progress) float64 {

	_, no := target.Dimensions()
	nj := uint(len(config.Importance))

	score, reject, refine := 0.0, true, false
	for i := uint(0); i < no; i++ {
		j := i % nj

		if config.Importance[j] == 0 {
			continue
		}

		s := math.Abs(location.Surplus[i])
		if s >= config.Rejection[j] {
			reject = false
		}
		if s > config.Refinement[j] {
			refine = true
		}

		score += config.Importance[j] * s
	}

	if reject {
		return -1
	}
	if !refine {
		return 0
	}

	return score
}
