package internal

import (
	"fmt"

	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"
	"github.com/ready-steady/spline"

	"../../pkg/cache"
)

type sliceTarget struct {
	problem *Problem
	config  *TargetConfig

	coreIndex []uint
	timeline  []float64

	cache *cache.Cache
}

func newSliceTarget(p *Problem, c *TargetConfig) (*sliceTarget, error) {
	const (
		cacheCapacity = 10000
	)

	// The cores of interest.
	coreIndex, err := enumerate(p.system.nc, c.CoreIndex)
	if err != nil {
		return nil, err
	}

	// The time moments of interest.
	timeline, err := subdivide(p.system.schedule.Span, c.TimeStep, c.TimeFraction)
	if err != nil {
		return nil, err
	}

	target := &sliceTarget{
		problem: p,
		config:  c,

		coreIndex: coreIndex,
		timeline:  timeline,

		cache: cache.New(cacheCapacity),
	}

	return target, nil
}

func (t *sliceTarget) String() string {
	return String(t)
}

func (t *sliceTarget) Config() *TargetConfig {
	return t.config
}

func (t *sliceTarget) Dimensions() (uint, uint) {
	return 1 + t.problem.nz, uint(len(t.coreIndex)) * 2 // +1 for time
}

func (t *sliceTarget) Compute(node, value []float64) {
	p := t.problem
	s := p.system

	var interpolant *spline.Cubic
	var key string

	key = stringify(node[1:]) // +1 for time
	if result, ok := t.cache.Get(key); ok {
		interpolant = result.(*spline.Cubic)
	}

	left, right := t.timeline[0], t.timeline[len(t.timeline)-1]

	if interpolant == nil {
		schedule := s.computeSchedule(p.transform(node[1:])) // +1 for time
		Q, time, err := s.temperature.Compute(s.power.Process(schedule), []float64{0, right})
		if err != nil {
			panic("cannot compute a temperature profile")
		}

		i, j := locate(left, right, time)

		time = time[i:j]
		Q = slice(Q[i*s.nc:j*s.nc], t.coreIndex, s.nc)

		interpolant = spline.NewCubic(time, Q)

		t.cache.Add(key, interpolant)
	}

	Q := interpolant.Evaluate([]float64{left + (right-left)*node[0]})

	for i := range Q {
		value[i*2] = Q[i]
		value[i*2+1] = Q[i] * Q[i]
	}
}

func (t *sliceTarget) Refine(node, surplus []float64, volume float64) float64 {
	return Refine(t, node, surplus, volume)
}

func (t *sliceTarget) Monitor(level, np, na uint) {
	if !t.config.Verbose {
		return
	}
	if level == 0 {
		fmt.Printf("%10s %15s %15s %15s\n",
			"Level", "Passive Nodes", "Active Nodes", "Evaluations")
	}
	fmt.Printf("%10d %15d %15d %15d\n", level, np, na, t.cache.Length())
}

func (t *sliceTarget) Generate(ns uint) []float64 {
	ni, _ := t.Dimensions()

	timeline := t.timeline
	nt := uint(len(timeline))

	left, right := timeline[0], timeline[nt-1]

	uniform := uniform.New(0, 1)
	points := make([]float64, ns*nt*ni)

	for i, k := uint(0), uint(0); i < ns; i++ {
		sample := probability.Sample(uniform, ni-1) // -1 for time
		for j := uint(0); j < nt; j++ {
			points[k] = (timeline[j] - left) / (right - left)
			copy(points[k+1:], sample)
			k += ni
		}
	}

	return points
}
