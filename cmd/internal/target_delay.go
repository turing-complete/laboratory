package internal

import (
	"fmt"
)

type delayTarget struct {
	problem *Problem
}

func newDelayTarget(p *Problem) Target {
	return &delayTarget{problem: p}
}

func (t *delayTarget) Inputs() uint {
	return t.problem.nz
}

func (t *delayTarget) Outputs() uint {
	return 1
}

func (t *delayTarget) String() string {
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", t.Inputs(), t.Outputs())
}

func (t *delayTarget) Evaluate(node []float64, value []float64, _ []uint64) {
	p := t.problem
	value[0] = p.time.Recompute(p.schedule, p.transform(node)).Span
}

func (t *delayTarget) Progress(level uint32, na, nt uint) {
	if level == 0 {
		fmt.Printf("%10s %15s %15s\n", "Level", "Passive Nodes", "Active Nodes")
	}
	fmt.Printf("%10d %15d %15d\n", level, nt-na, na)
}
