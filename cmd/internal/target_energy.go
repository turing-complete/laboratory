package internal

import (
	"fmt"
)

type energyTarget struct {
	problem *Problem
}

func newEnergyTarget(p *Problem) Target {
	return &energyTarget{problem: p}
}

func (t *energyTarget) InputsOutputs() (uint32, uint32) {
	return t.problem.zc, 1
}

func (t *energyTarget) String() string {
	ic, oc := t.InputsOutputs()
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", ic, oc)
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
}

func (t *energyTarget) Progress(level uint8, activeNodes, totalNodes uint32) {
	if level == 0 {
		t.problem.Printf("%10s %15s %15s\n",
			"Level", "Passive Nodes", "Active Nodes")
	}

	passiveNodes := totalNodes - activeNodes
	t.problem.Printf("%10d %15d %15d\n", level, passiveNodes, activeNodes)
}
