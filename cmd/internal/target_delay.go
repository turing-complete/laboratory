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

func (t *delayTarget) InputsOutputs() (uint32, uint32) {
	return t.problem.zc, 1
}

func (t *delayTarget) String() string {
	ic, oc := t.InputsOutputs()
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", ic, oc)
}

func (t *delayTarget) Evaluate(node []float64, value []float64, _ []uint64) {
	p := t.problem
	value[0] = p.time.Recompute(p.schedule, p.transform(node)).Span
}

func (t *delayTarget) Progress(level uint8, activeNodes, totalNodes uint32) {
	if level == 0 {
		t.problem.Printf("%10s %15s %15s\n",
			"Level", "Passive Nodes", "Active Nodes")
	}

	passiveNodes := totalNodes - activeNodes
	t.problem.Printf("%10d %15d %15d\n", level, passiveNodes, activeNodes)
}
