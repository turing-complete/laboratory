package internal

import (
	"errors"
	"fmt"

	"github.com/ready-steady/simulation/power"
	"github.com/ready-steady/simulation/temperature/numeric"
)

type temperatureTarget struct {
	problem *Problem
	config  *TargetConfig

	power       *power.Power
	temperature *numeric.Temperature

	cores    []uint
	timeline []float64
	shift    bool
}

func newTemperatureTarget(p *Problem, tac *TargetConfig,
	tec *TemperatureConfig) (*temperatureTarget, error) {

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

	Δt := tac.TimeStep
	if Δt <= 0 {
		return nil, errors.New("the time step should be positive")
	}

	// The time moments of interest.
	interval := tac.TimeInterval
	switch len(interval) {
	case 0:
		interval = []float64{0, p.schedule.Span}
	case 1:
		interval = []float64{interval[0], interval[0]}
	default:
	}
	if interval[0] < 0 || interval[0] > interval[1] || interval[1] > p.schedule.Span {
		return nil, errors.New(fmt.Sprintf(
			"the time interval should be between 0 and %g seconds", p.schedule.Span))
	}

	timeline := []float64{}
	for t := interval[0]; t <= interval[1]; t += Δt {
		timeline = append(timeline, t)
	}

	// Force the first time moment to be zero.
	shift := timeline[0] != 0
	if shift {
		timeline = append([]float64{0}, timeline...)
	}

	target := &temperatureTarget{
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

func (t *temperatureTarget) String() string {
	ni, no := t.Dimensions()
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", ni, no)
}

func (t *temperatureTarget) Dimensions() (uint, uint) {
	nci, ns := uint(len(t.cores)), uint(len(t.timeline))
	if t.shift {
		ns--
	}
	return t.problem.nz, ns * nci * 2
}

func (t *temperatureTarget) Compute(node, value []float64) {
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

func (t *temperatureTarget) Refine(surplus []float64) bool {
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

func (t *temperatureTarget) Monitor(level, np, na uint) {
	if !t.config.Verbose {
		return
	}
	if level == 0 {
		fmt.Printf("%10s %15s %15s\n", "Level", "Passive Nodes", "Active Nodes")
	}
	fmt.Printf("%10d %15d %15d\n", level, np, na)
}
