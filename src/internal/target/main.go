package target

import (
	"errors"
	"runtime"
	"sync"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

var (
	Workers = uint(runtime.GOMAXPROCS(0))
)

type Target interface {
	Dimensions() (uint, uint)
	Compute([]float64, []float64)
	Forward([]float64) []float64
	Inverse([]float64) []float64
}

func New(system *system.System, uncertainty *uncertainty.Uncertainty,
	config *config.Target) (Target, error) {

	switch config.Name {
	case "end-to-end-delay":
		return newDelay(system, uncertainty, config)
	case "total-energy":
		return newEnergy(system, uncertainty, config)
	case "maximal-temperature":
		return newTemperature(system, uncertainty, config)
	default:
		return nil, errors.New("the target is unknown")
	}
}

func Invoke(target Target, points []float64) []float64 {
	ni, no := target.Dimensions()
	np := uint(len(points)) / ni

	values := make([]float64, np*no)

	jobs := make(chan uint, np)
	group := sync.WaitGroup{}
	group.Add(int(np))

	for i := uint(0); i < Workers; i++ {
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
