package internal

import (
	"errors"
	"fmt"

	"github.com/ready-steady/simulation/power"
	"github.com/ready-steady/simulation/temperature/numeric"
)

type temperatureTarget struct {
	problem *Problem

	cores    []uint
	Δt       float64
	shift    bool
	timeline []float64

	power       *power.Power
	temperature *numeric.Temperature
}

func newTemperatureTarget(p *Problem) (Target, error) {
	c := &p.Config.Temperature

	power := power.New(p.platform, p.application)
	temperature, err := numeric.New(&c.Config)
	if err != nil {
		return nil, err
	}

	// The cores of interest.
	cores := c.CoreIndex
	if len(cores) == 0 {
		cores = make([]uint, p.nc)
		for i := uint(0); i < p.nc; i++ {
			cores[i] = i
		}
	}

	Δt := c.TimeStep
	if Δt <= 0 {
		return nil, errors.New("the time step should be positive")
	}

	// The time moments of interest.
	interval := c.TimeInterval
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

		cores:    cores,
		Δt:       Δt,
		timeline: timeline,
		shift:    shift,

		power:       power,
		temperature: temperature,
	}

	return target, nil
}

func (t *temperatureTarget) Inputs() uint {
	return t.problem.nz
}

func (t *temperatureTarget) Outputs() uint {
	nci, ns := uint(len(t.cores)), uint(len(t.timeline))
	if t.shift {
		ns--
	}
	return ns * nci
}

func (t *temperatureTarget) String() string {
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", t.Inputs(), t.Outputs())
}

func (t *temperatureTarget) Evaluate(node, value []float64, _ []uint64) {
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

	for i := uint(0); i < ns; i++ {
		for j := uint(0); j < nci; j++ {
			value[i*nci+j] = Q[i*nc+cores[j]]
		}
	}
}

func (t *temperatureTarget) Progress(level uint32, na, nt uint) {
	if level == 0 {
		fmt.Printf("%10s %15s %15s\n", "Level", "Passive Nodes", "Active Nodes")
	}
	fmt.Printf("%10d %15d %15d\n", level, nt-na, na)
}
