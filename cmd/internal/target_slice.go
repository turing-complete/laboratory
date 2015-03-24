package internal

import (
	"fmt"

	"github.com/ready-steady/numeric/interpolation/spline"
	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"
	"github.com/ready-steady/simulation/power"
	"github.com/ready-steady/simulation/temperature/numeric"

	"../../pkg/cache"
)

type sliceTarget struct {
	problem *Problem
	config  *TargetConfig

	power       *power.Power
	temperature *numeric.Temperature

	cores    []uint
	timeline []float64

	cache *cache.Cache
}

func newSliceTarget(p *Problem, tac *TargetConfig,
	tec *TemperatureConfig) (*sliceTarget, error) {

	const (
		cacheCapacity = 10000
	)

	power := power.New(p.platform, p.application)
	temperature, err := numeric.New(&tec.Config)
	if err != nil {
		return nil, err
	}

	// The cores of interest.
	cores := tac.CoreIndex
	if len(cores) == 0 {
		cores = make([]uint, p.nc)
		for i := uint(0); i < p.nc; i++ {
			cores[i] = i
		}
	}

	// The time moments of interest.
	timeline, err := subdivide(p.schedule.Span, tac.TimeStep, tac.TimeFraction)
	if err != nil {
		return nil, err
	}

	target := &sliceTarget{
		problem: p,
		config:  tac,

		power:       power,
		temperature: temperature,

		cores:    cores,
		timeline: timeline,

		cache: cache.New(cacheCapacity),
	}

	return target, nil
}

func (t *sliceTarget) String() string {
	return CommonTarget{t}.String()
}

func (t *sliceTarget) Config() *TargetConfig {
	return t.config
}

func (t *sliceTarget) Dimensions() (uint, uint) {
	return 1 + t.problem.nz, uint(len(t.cores)) * 2 // +1 for time
}

func (t *sliceTarget) Compute(node, value []float64) {
	p := t.problem

	var interpolant *spline.Cubic
	var key string

	key = stringify(node[1:]) // +1 for time
	if result, ok := t.cache.Get(key); ok {
		interpolant = result.(*spline.Cubic)
	}

	left, right := t.timeline[0], t.timeline[len(t.timeline)-1]

	if interpolant == nil {
		schedule := p.time.Recompute(p.schedule, p.transform(node[1:])) // +1 for time

		Q, time, err := t.temperature.Compute(t.power.Process(schedule), []float64{0, right})
		if err != nil {
			panic("cannot compute a temperature profile")
		}

		i, j := locate(left, right, time)

		time = time[i:j]
		Q = slice(Q[i*p.nc:j*p.nc], t.cores, p.nc)

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
	return CommonTarget{t}.Refine(node, surplus, volume)
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
