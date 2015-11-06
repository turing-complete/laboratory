package target

import (
	"errors"
	"sync"

	"github.com/ready-steady/adapt"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

type Target adapt.Target

func New(system *system.System, config *config.Target) (Target, error) {
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
