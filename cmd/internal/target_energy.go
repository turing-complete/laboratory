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

func (t *energyTarget) Inputs() uint {
	return t.problem.zc
}

func (t *energyTarget) Outputs() uint {
	return 1
}

func (t *energyTarget) Pseudos() uint {
	return 0
}

func (t *energyTarget) String() string {
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", t.Inputs(), t.Outputs())
}

func (t *energyTarget) Evaluate(node, value []float64, _ []uint64) {
	p := t.problem

	cores, tasks := p.platform.Cores, p.application.Tasks
	schedule := p.time.Recompute(p.schedule, p.transform(node))

	value[0] = 0
	for i := range tasks {
		value[0] += (schedule.Finish[i] - schedule.Start[i]) *
			cores[schedule.Mapping[i]].Power[tasks[i].Type]
	}
}

func (t *energyTarget) Progress(level uint32, activeNodes, totalNodes uint) {
	if level == 0 {
		fmt.Printf("%10s %15s %15s\n", "Level", "Passive Nodes", "Active Nodes")
	}

	passiveNodes := totalNodes - activeNodes
	fmt.Printf("%10d %15d %15d\n", level, passiveNodes, activeNodes)
}
