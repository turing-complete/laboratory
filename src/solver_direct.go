package main

import (
	"fmt"

	"github.com/ready-steady/numan/interp/adhier"
)

type directSolver struct {
	*baseSolver
}

func (s *directSolver) Construct() *adhier.Surrogate {
	c := &s.problem.config

	ic, oc := s.ic, s.oc
	NC := uint32(0)

	jobs := s.spawnWorkers()

	if c.Verbose {
		fmt.Printf("%12s %12s\n", "new nodes", "total nodes")
	}

	surrogate := s.interpolator.Compute(func(nodes []float64, index []uint64) []float64 {
		nc := uint32(len(nodes)) / ic
		NC += nc

		if c.Verbose {
			fmt.Printf("%12d %12d\n", nc, NC)
		}

		done := make(chan result, nc)
		values := make([]float64, oc*nc)

		for i := uint32(0); i < nc; i++ {
			jobs <- job{
				node:  nodes[i*ic:],
				value: values[i*oc:],
				done:  done,
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
