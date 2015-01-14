package solver

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/ready-steady/numan/basis/linhat"
	"github.com/ready-steady/numan/grid/newcot"
	"github.com/ready-steady/numan/interp/adhier"
)

type Job struct {
	Key   string
	Data  []float64
	Node  []float64
	Value []float64
	Done  chan<- Result
}

type Result struct {
	Key  string
	Data []float64
}

type Config struct {
	Inputs  uint16 // The number of inputs.
	Outputs uint16 // The number of outputs.

	// The number specifying how many of the inputs should be used for caching.
	CacheInputs uint16

	// The number of workers evaluating of the quantity of interest.
	Workers uint8
	// The configuration of the algorithm for interpolation.
	Interpolation adhier.Config

	Verbose bool // A flag for displaying progress information.
}

type Solver struct {
	config       Config
	target       func(<-chan Job)
	interpolator *adhier.Interpolator
}

func New(config Config, target func(<-chan Job)) (*Solver, error) {
	if config.Interpolation.AbsError <= 0 {
		return nil, errors.New("the absolute-error tolerance is invalid")
	}
	if config.Interpolation.RelError <= 0 {
		return nil, errors.New("the relative-error tolerance is invalid")
	}

	interpolator := adhier.New(newcot.NewOpen(config.Inputs),
		linhat.NewOpen(config.Inputs), adhier.Config(config.Interpolation),
		config.Outputs)

	solver := &Solver{
		config:       config,
		target:       target,
		interpolator: interpolator,
	}

	return solver, nil
}

func (s *Solver) Construct() *adhier.Surrogate {
	if s.config.CacheInputs == 0 {
		return s.constructDirect()
	} else {
		return s.constructCached()
	}
}

func (s *Solver) Compute(nodes []float64) []float64 {
	ic, oc := uint32(s.config.Inputs), uint32(s.config.Outputs)
	nc := uint32(len(nodes)) / ic

	jobs := s.spawnWorkers()

	done := make(chan Result, nc)
	values := make([]float64, oc*nc)

	jc, rc := uint32(0), uint32(0)
	nextJob := Job{
		Node:  nodes[jc*ic:],
		Value: values[jc*oc:],
		Done:  done,
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

			nextJob = Job{
				Node:  nodes[jc*ic:],
				Value: values[jc*oc:],
				Done:  done,
			}
		case <-done:
			rc++
		}
	}

	return values
}

func (s *Solver) Evaluate(surrogate *adhier.Surrogate, points []float64) []float64 {
	return s.interpolator.Evaluate(surrogate, points)
}

func (s *Solver) spawnWorkers() chan<- Job {
	wc := int(s.config.Workers)
	if wc <= 0 {
		wc = runtime.NumCPU()
	}

	if s.config.Verbose {
		fmt.Printf("Using %d workers...\n", wc)
	}

	runtime.GOMAXPROCS(wc)

	jobs := make(chan Job)
	for i := 0; i < wc; i++ {
		go s.target(jobs)
	}

	return jobs
}
