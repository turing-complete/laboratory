package internal

import (
	"github.com/ready-steady/simulation/power"
	"github.com/ready-steady/simulation/temperature/numeric"
)

type profileTarget struct {
	problem *Problem
	config  *TargetConfig

	power       *power.Power
	temperature *numeric.Temperature

	cores    []uint
	timeline []float64
	shift    uint
}

func newProfileTarget(p *Problem, tac *TargetConfig,
	tec *TemperatureConfig) (*profileTarget, error) {

	power := power.New(p.platform, p.application)
	temperature, err := numeric.New(&tec.Config)
	if err != nil {
		return nil, err
	}

	// The cores of interest.
	cores, err := enumerate(p.nc, tac.CoreIndex)
	if err != nil {
		return nil, err
	}

	// The time moments of interest.
	timeline, err := subdivide(p.schedule.Span, tac.TimeStep, tac.TimeFraction)
	if err != nil {
		return nil, err
	}

	shift := uint(0)

	// Force the first time moment to be zero.
	if timeline[0] != 0 {
		shift++
		timeline = append([]float64{0}, timeline...)
	}

	// Make sure to have at least three time moments.
	if len(timeline) == 2 {
		shift++
		timeline = []float64{0, timeline[1] / 2, timeline[1]}
	}

	target := &profileTarget{
		problem: p,
		config:  tac,

		power:       power,
		temperature: temperature,

		cores:    cores,
		timeline: timeline,
		shift:    shift,
	}

	return target, nil
}

func (t *profileTarget) String() string {
	return CommonTarget{t}.String()
}

func (t *profileTarget) Dimensions() (uint, uint) {
	nci, ns := uint(len(t.cores)), uint(len(t.timeline))-t.shift
	return t.problem.nz, ns * nci * 2
}

func (t *profileTarget) Compute(node, value []float64) {
	p := t.problem

	schedule := p.time.Recompute(p.schedule, p.transform(node))
	Q, _, err := t.temperature.Compute(t.power.Process(schedule), t.timeline)
	if err != nil {
		panic("cannot compute a temperature profile")
	}

	cores := t.cores
	nc, nci, ns := p.nc, uint(len(cores)), uint(len(t.timeline))-t.shift

	Q = Q[t.shift*nc:]

	for i, k := uint(0), uint(0); i < ns; i++ {
		for j := uint(0); j < nci; j++ {
			value[k] = Q[i*nc+cores[j]]
			value[k+1] = value[k] * value[k]
			k += 2
		}
	}
}

func (t *profileTarget) Refine(surplus []float64, dimensions []bool) {
	nm, ε := uint(len(surplus))/2, t.config.Tolerance

	refine := false

	for i := uint(0); i < nm; i++ {
		if surplus[i*2] > ε || -surplus[i*2] > ε {
			refine = true
			break
		}
	}

	for i := range dimensions {
		dimensions[i] = refine
	}
}

func (t *profileTarget) Monitor(level, np, na uint) {
	if t.config.Verbose {
		CommonTarget{t}.Monitor(level, np, na)
	}
}

func (t *profileTarget) Generate(ns uint) []float64 {
	return CommonTarget{t}.Generate(ns)
}
