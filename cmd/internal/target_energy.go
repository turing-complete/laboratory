package internal

import (
	"fmt"
)

type energyTarget struct {
	problem *Problem
	config  *TargetConfig
}

func newEnergyTarget(p *Problem, c *TargetConfig) *energyTarget {
	return &energyTarget{
		problem: p,
		config:  c,
	}
}

func (t *energyTarget) String() string {
	ni, no := t.Dimensions()
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", ni, no)
}

func (t *energyTarget) Dimensions() (uint, uint) {
	return t.problem.nz, 1
}

func (t *energyTarget) Compute(node, value []float64) {
	p := t.problem

	cores, tasks := p.platform.Cores, p.application.Tasks
	schedule := p.time.Recompute(p.schedule, p.transform(node))

	value[0] = 0
	for i := range tasks {
		value[0] += (schedule.Finish[i] - schedule.Start[i]) *
			cores[schedule.Mapping[i]].Power[tasks[i].Type]
	}
}

func (t *energyTarget) Refine(surplus []float64) bool {
	return surplus[0] > t.config.Tolerance || -surplus[0] > t.config.Tolerance
}

func (t *energyTarget) Monitor(level, np, na uint) {
	if !t.config.Verbose {
		return
	}
	if level == 0 {
		fmt.Printf("%10s %15s %15s\n", "Level", "Passive Nodes", "Active Nodes")
	}
	fmt.Printf("%10d %15d %15d\n", level, np, na)
}
