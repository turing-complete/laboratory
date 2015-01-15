package solver

import (
	"fmt"

	"github.com/ready-steady/numeric/interpolation/adhier"

	"../cache"
)

const (
	cacheCapacity = 1000
)

func (s *Solver) constructCached() *adhier.Surrogate {
	verbose := s.config.Verbose

	ic, oc, cc := uint32(s.config.Inputs), uint32(s.config.Outputs), uint32(s.config.CacheInputs)
	NC, EC := uint32(0), uint32(0)

	cache := cache.New(ic-cc, cacheCapacity)
	jobs := s.spawnWorkers()

	if verbose {
		fmt.Printf("%12s %12s (%6s) %12s %12s (%6s)\n",
			"New nodes", "New evals", "%", "Total nodes", "Total evals", "%")
	}

	surrogate := s.interpolator.Compute(func(nodes []float64, index []uint64) []float64 {
		nc, ec := uint32(len(nodes))/ic, uint32(0)

		NC += nc
		if verbose {
			fmt.Printf("%12d", nc)
		}

		done := make(chan Result, nc)
		values := make([]float64, oc*nc)

		for i := uint32(0); i < nc; i++ {
			key := cache.Key(index[cc+i*ic:])

			data := cache.Get(key)
			if data == nil {
				ec++
			}

			jobs <- Job{
				Key:   key,
				Data:  data,
				Node:  nodes[i*ic:],
				Value: values[i*oc:],
				Done:  done,
			}
		}

		for i := uint32(0); i < nc; i++ {
			result := <-done
			cache.Set(result.Key, result.Data)
		}

		EC += ec
		if verbose {
			fmt.Printf(" %12d (%6.2f) %12d %12d (%6.2f)\n",
				ec, float64(ec)/float64(nc)*100,
				NC, EC, float64(EC)/float64(NC)*100)
		}

		return values
	})

	close(jobs)

	return surrogate
}
