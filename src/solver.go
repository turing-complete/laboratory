package main

import (
	"fmt"
	"runtime"

	"github.com/ready-steady/numan/basis/linhat"
	"github.com/ready-steady/numan/grid/newcot"
	"github.com/ready-steady/numan/interp/adhier"
)

type solver struct {
	problem *problem
	target  target

	ic uint32 // inputs
	oc uint32 // outputs
	fc uint32 // fake inputs like time

	interpolator *adhier.Self
}

func newSolver(problem *problem, target target) *solver {
	ic, oc := target.InputsOutputs()

	interpolator := adhier.New(newcot.NewOpen(uint16(ic)), linhat.NewOpen(uint16(ic)),
		adhier.Config(problem.config.Interpolation), uint16(oc))

	return &solver{
		problem: problem,
		target:  target,

		ic: ic,
		oc: oc,
		fc: ic - problem.zc,

		interpolator: interpolator,
	}
}

func (s *solver) hasFakeInputs() bool {
	return s.fc > 0
}

func (s *solver) construct() *adhier.Surrogate {
	if s.hasFakeInputs() {
		// NOTE: Caching only make sense when there are fake inputs.
		return s.constructCached()
	} else {
		return s.constructDirect()
	}
}

func (s *solver) compute(nodes []float64) []float64 {
	ic, oc := s.ic, s.oc
	nc := uint32(len(nodes)) / ic

	jobs := s.spawnWorkers()

	done := make(chan result, nc)
	values := make([]float64, oc*nc)

	jc, rc := uint32(0), uint32(0)
	nextJob := job{
		node:  nodes[jc*ic:],
		value: values[jc*oc:],
		done:  done,
	}

	for jc < nc || rc < nc {
		select {
		case jobs <- nextJob:
			jc++

			if jc >= nc {
				close(jobs)
				jobs = nil
				continue
			}

			nextJob = job{
				node:  nodes[jc*ic:],
				value: values[jc*oc:],
				done:  done,
			}
		case <-done:
			rc++
		}
	}

	return values
}

func (s *solver) evaluate(surrogate *adhier.Surrogate, points []float64) []float64 {
	return s.interpolator.Evaluate(surrogate, points)
}

func (s *solver) spawnWorkers() chan<- job {
	c := &s.problem.config

	wc := int(c.Workers)
	if wc <= 0 {
		wc = runtime.NumCPU()
	}

	if c.Verbose {
		fmt.Printf("Using %d workers...\n", wc)
	}

	runtime.GOMAXPROCS(wc)

	jobs := make(chan job)
	for i := 0; i < wc; i++ {
		go s.target.Serve(jobs)
	}

	return jobs
}
