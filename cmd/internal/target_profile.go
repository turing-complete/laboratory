package internal

import (
	"errors"
	"fmt"

	"github.com/ready-steady/simulation/power"
	"github.com/ready-steady/simulation/temperature/numeric"
)

type profileTarget struct {
	problem *Problem

	Δt       float64
	shift    bool
	timeline []float64

	power       *power.Power
	temperature *numeric.Temperature
}

func newProfileTarget(p *Problem) (Target, error) {
	c := &p.Config

	power := power.New(p.platform, p.application)
	temperature, err := numeric.New((*numeric.Config)(&c.TempAnalysis))
	if err != nil {
		return nil, err
	}

	Δt := c.TempAnalysis.TimeStep

	stepIndex := c.StepIndex
	sc := uint(len(stepIndex))

	if sc == 0 {
		sc = uint(p.schedule.Span / Δt)
		stepIndex = make([]uint, sc)
		for i := uint(0); i < sc; i++ {
			stepIndex[i] = i
		}
	}

	// Force the first index to be zero.
	shift := stepIndex[0] != 0
	if shift {
		newIndex := make([]uint, sc+1)
		copy(newIndex[1:], stepIndex)
		stepIndex = newIndex
		sc++
	}

	timeline := make([]float64, sc)
	for i, max := uint(0), uint(p.schedule.Span/Δt)-1; i < sc; i++ {
		if stepIndex[i] > max {
			return nil, errors.New(fmt.Sprintf("the step indices should not exceed %d", max))
		}
		timeline[i] = float64(stepIndex[i]) * Δt
	}

	target := &profileTarget{
		problem: p,

		Δt:       Δt,
		shift:    shift,
		timeline: timeline,

		power:       power,
		temperature: temperature,
	}

	return target, nil
}

func (t *profileTarget) Inputs() uint {
	return t.problem.zc
}

func (t *profileTarget) Outputs() uint {
	cci, sc := uint(len(t.problem.Config.CoreIndex)), uint(len(t.timeline))
	if t.shift {
		sc--
	}
	return sc * cci
}

func (t *profileTarget) Pseudos() uint {
	return 0
}

func (t *profileTarget) String() string {
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", t.Inputs(), t.Outputs())
}

func (t *profileTarget) Evaluate(node, value []float64, _ []uint64) {
	p := t.problem

	schedule := p.time.Recompute(p.schedule, p.transform(node))
	Q, _, err := t.temperature.Compute(t.power.Process(schedule), t.timeline)
	if err != nil {
		panic("cannot compute a temperature profile")
	}

	coreIndex := p.Config.CoreIndex
	cc, cci, sc := p.cc, uint(len(coreIndex)), uint(len(t.timeline))

	if t.shift {
		Q = Q[cc:]
		sc--
	}

	for i := uint(0); i < sc; i++ {
		for j := uint(0); j < cci; j++ {
			value[i*cci+j] = Q[i*cc+coreIndex[j]]
		}
	}
}

func (t *profileTarget) Progress(level uint32, activeNodes, totalNodes uint) {
	if level == 0 {
		fmt.Printf("%10s %15s %15s\n",
			"Level", "Passive Nodes", "Active Nodes")
	}

	passiveNodes := totalNodes - activeNodes
	fmt.Printf("%10d %15d %15d\n", level, passiveNodes, activeNodes)
}
