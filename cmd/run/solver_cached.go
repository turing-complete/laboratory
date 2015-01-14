package main

import (
	"fmt"

	"github.com/ready-steady/numan/interp/adhier"

	"../../pkg/cache"
)

const (
	cacheCapacity = 1000
)

type cachedSolver struct {
	*baseSolver

	fc uint32 // fake inputs like time
}

func (s *cachedSolver) Construct() *adhier.Surrogate {
	p := s.problem
	c := &p.config

	ic, oc, fc := s.ic, s.oc, s.fc
	NC, EC := uint32(0), uint32(0)

	cache := cache.New(p.zc, cacheCapacity)
	jobs := s.spawnWorkers()

	if c.Verbose {
		fmt.Printf("%12s %12s (%6s) %12s %12s (%6s)\n",
			"New nodes", "New evals", "%", "Total nodes", "Total evals", "%")
	}

	surrogate := s.interpolator.Compute(func(nodes []float64, index []uint64) []float64 {
		nc, ec := uint32(len(nodes))/ic, uint32(0)

		NC += nc
		if c.Verbose {
			fmt.Printf("%12d", nc)
		}

		done := make(chan result, nc)
		values := make([]float64, oc*nc)

		for i := uint32(0); i < nc; i++ {
			key := cache.Key(index[fc+i*ic:])

			data := cache.Get(key)
			if data == nil {
				ec++
			}

			jobs <- job{
				key:   key,
				data:  data,
				node:  nodes[i*ic:],
				value: values[i*oc:],
				done:  done,
			}
		}

		for i := uint32(0); i < nc; i++ {
			result := <-done
			cache.Set(result.key, result.data)
		}

		EC += ec
		if c.Verbose {
			fmt.Printf(" %12d (%6.2f) %12d %12d (%6.2f)\n",
				ec, float64(ec)/float64(nc)*100,
				NC, EC, float64(EC)/float64(NC)*100)
		}

		return values
	})

	close(jobs)

	return surrogate
}
