package internal

import (
	"fmt"
)

type delayTarget struct {
	problem *Problem
	config  *TargetConfig
}

func newDelayTarget(p *Problem, c *TargetConfig) *delayTarget {
	return &delayTarget{
		problem: p,
		config:  c,
	}
}

func (t *delayTarget) String() string {
	ni, no := t.Dimensions()
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", ni, no)
}

func (t *delayTarget) Dimensions() (uint, uint) {
	return t.problem.nz, 1
}

func (t *delayTarget) Compute(node []float64, value []float64) {
	p := t.problem
	value[0] = p.time.Recompute(p.schedule, p.transform(node)).Span
}

func (t *delayTarget) Refine(surplus []float64) bool {
	return surplus[0] > t.config.Tolerance || -surplus[0] > t.config.Tolerance
}

func (t *delayTarget) Monitor(level, np, na uint) {
	if !t.config.Verbose {
		return
	}
	if level == 0 {
		fmt.Printf("%10s %15s %15s\n", "Level", "Passive Nodes", "Active Nodes")
	}
	fmt.Printf("%10d %15d %15d\n", level, np, na)
}
