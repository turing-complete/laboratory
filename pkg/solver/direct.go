package solver

import (
	"fmt"

	"github.com/ready-steady/numeric/interpolation/adhier"
)

func (s *Solver) constructDirect() *adhier.Surrogate {
	verbose := s.config.Verbose

	ic, oc := uint32(s.config.Inputs), uint32(s.config.Outputs)
	NC := uint32(0)

	jobs := s.spawnWorkers()

	if verbose {
		fmt.Printf("%12s %12s\n", "New nodes", "Total nodes")
	}

	surrogate := s.interpolator.Compute(func(nodes []float64, _ []uint64) []float64 {
		nc := uint32(len(nodes)) / ic
		NC += nc

		if verbose {
			fmt.Printf("%12d %12d\n", nc, NC)
		}

		done := make(chan Result, nc)
		values := make([]float64, nc*oc)

		for i := uint32(0); i < nc; i++ {
			jobs <- Job{
				Node:  nodes[i*ic : (i+1)*ic],
				Value: values[i*oc : (i+1)*oc],
				Done:  done,
			}
		}

		for i := uint32(0); i < nc; i++ {
			<-done
		}

		return values
	})

	close(jobs)

	return surrogate
}
