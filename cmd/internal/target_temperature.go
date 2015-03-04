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
	cores := c.Cores
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

	steps := c.Steps
	ns := uint(len(steps))

	if ns == 0 {
		ns = uint(p.schedule.Span / Δt)
		steps = make([]uint, ns)
		for i := uint(0); i < ns; i++ {
			steps[i] = i
		}
	}

	// Force the first index to be zero.
	shift := steps[0] != 0
	if shift {
		newSteps := make([]uint, ns+1)
		copy(newSteps[1:], steps)
		steps = newSteps
		ns++
	}

	timeline := make([]float64, ns)
	for i, max := uint(0), uint(p.schedule.Span/Δt)-1; i < ns; i++ {
		if steps[i] > max {
			return nil, errors.New(fmt.Sprintf("the step indices should not exceed %d", max))
		}
		timeline[i] = float64(steps[i]) * Δt
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
