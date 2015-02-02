package internal

import (
	"fmt"
	"sync/atomic"
)

type energyTarget struct {
	problem *Problem

	ec uint32
}

func newEnergyTarget(p *Problem) Target {
	return &energyTarget{problem: p}
}

func (t *energyTarget) Evaluate(node, value []float64, _ []uint64) {
	p := t.problem

	cores, tasks := p.platform.Cores, p.application.Tasks
	schedule := p.time.Recompute(p.schedule, p.transform(node))

	value[0] = 0
	for i := range tasks {
		value[0] += (schedule.Finish[i] - schedule.Start[i]) *
			cores[uint32(schedule.Mapping[i])].Power[tasks[i].Type]
	}

	atomic.AddUint32(&t.ec, 1)
}

func (t *energyTarget) InputsOutputs() (uint32, uint32) {
	return t.problem.zc, 1
}

func (t *energyTarget) Evaluations() uint32 {
	return t.ec
}

func (t *energyTarget) String() string {
	ic, oc := t.InputsOutputs()
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", ic, oc)
}
