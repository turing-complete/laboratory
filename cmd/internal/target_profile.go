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
	shift    bool
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
	timeline, err := subdivide(tac.TimeInterval, tac.TimeStep, p.schedule.Span)
	if err != nil {
		return nil, err
	}

	// Force the first time moment to be zero.
	shift := timeline[0] != 0
	if shift {
		timeline = append([]float64{0}, timeline...)
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
	nci, ns := uint(len(t.cores)), uint(len(t.timeline))
	if t.shift {
		ns--
	}
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
	nc, nci, ns := p.nc, uint(len(cores)), uint(len(t.timeline))

	if t.shift {
		Q = Q[nc:]
		ns--
	}

	for i, k := uint(0), uint(0); i < ns; i++ {
		for j := uint(0); j < nci; j++ {
			value[k] = Q[i*nc+cores[j]]
			value[k+1] = value[k] * value[k]
			k += 2
		}
	}
}

func (t *profileTarget) Refine(surplus []float64) bool {
	no, nci := uint(len(surplus)), uint(len(t.cores))
	ε := t.config.Tolerance

	// The beginning.
	for i := uint(0); i < nci; i++ {
		if surplus[i*2] > ε || -surplus[i*2] > ε {
			return true
		}
	}

	surplus = surplus[no-nci*2:]

	// The ending.
	for i := uint(0); i < nci; i++ {
		if surplus[i*2] > ε || -surplus[i*2] > ε {
			return true
		}
	}

	return false
}

func (t *profileTarget) Monitor(level, np, na uint) {
	if t.config.Verbose {
		CommonTarget{t}.Monitor(level, np, na)
	}
}

func (t *profileTarget) Generate(ns uint) []float64 {
	return CommonTarget{t}.Generate(ns)
}
