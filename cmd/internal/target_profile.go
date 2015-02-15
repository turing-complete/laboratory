package internal

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ready-steady/simulation/power"
	"github.com/ready-steady/simulation/temperature"
)

type profileTarget struct {
	problem *Problem

	sc        uint
	stepIndex []uint

	power       *power.Power
	temperature *temperature.Temperature

	pool *sync.Pool
}

type profileData struct {
	P []float64
	S []float64
	Q []float64
}

func newProfileTarget(p *Problem) (Target, error) {
	c := &p.Config

	power, err := power.New(p.platform, p.application, c.TempAnalysis.TimeStep)
	if err != nil {
		return nil, err
	}

	temperature, err := temperature.New(temperature.Config(c.TempAnalysis))
	if err != nil {
		return nil, err
	}

	cc, sc := p.cc, uint(p.schedule.Span/c.TempAnalysis.TimeStep)
	nc := temperature.Nodes

	var stepIndex []uint
	if len(c.StepIndex) == 0 {
		stepIndex = make([]uint, sc)
		for i := uint(0); i < sc; i++ {
			stepIndex[i] = i
		}
	} else {
		stepIndex = make([]uint, len(c.StepIndex))
		max := uint(0)
		for i := range stepIndex {
			stepIndex[i] = c.StepIndex[i]
			if stepIndex[i] >= sc {
				return nil, errors.New("the step index is invalid")
			}
			if stepIndex[i] > max {
				max = stepIndex[i]
			}
		}
		sc = max + 1
	}

	target := &profileTarget{
		problem: p,

		sc:        sc,
		stepIndex: stepIndex,

		power:       power,
		temperature: temperature,

		pool: &sync.Pool{
			New: func() interface{} {
				return &profileData{
					P: make([]float64, cc*sc),
					S: make([]float64, nc*sc),
					Q: make([]float64, cc*sc),
				}
			},
		},
	}

	return target, nil
}

func (t *profileTarget) Inputs() uint {
	return t.problem.zc
}

func (t *profileTarget) Outputs() uint {
	return uint(len(t.stepIndex) * len(t.problem.Config.CoreIndex))
}

func (t *profileTarget) Pseudos() uint {
	return 0
}

func (t *profileTarget) String() string {
	return fmt.Sprintf("Target{inputs: %d, outputs: %d, steps: %d}",
		t.Inputs(), t.Outputs(), t.sc)
}

func (t *profileTarget) Evaluate(node, value []float64, _ []uint64) {
	p := t.problem

	coreIndex, stepIndex := p.Config.CoreIndex, t.stepIndex

	cc, cci, sc, sci := p.cc, uint(len(coreIndex)), t.sc, uint(len(stepIndex))

	data := t.pool.Get().(*profileData)

	t.power.Compute(p.time.Recompute(p.schedule, p.transform(node)), data.P, sc)
	t.temperature.ComputeTransient(data.P, data.Q, data.S, sc)

	for i := uint(0); i < sci; i++ {
		for j := uint(0); j < cci; j++ {
			value[i*cci+j] = data.Q[stepIndex[i]*cc+coreIndex[j]]
		}
	}

	t.pool.Put(data)
}

func (t *profileTarget) Progress(level uint32, activeNodes, totalNodes uint) {
	if level == 0 {
		fmt.Printf("%10s %15s %15s\n",
			"Level", "Passive Nodes", "Active Nodes")
	}

	passiveNodes := totalNodes - activeNodes
	fmt.Printf("%10d %15d %15d\n", level, passiveNodes, activeNodes)
}
