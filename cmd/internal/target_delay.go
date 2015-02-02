package internal

import (
	"fmt"
	"sync/atomic"
)

type delayTarget struct {
	problem *Problem

	ec uint32
}

func newDelayTarget(p *Problem) Target {
	return &delayTarget{problem: p}
}

func (t *delayTarget) Evaluate(node []float64, value []float64, _ []uint64) {
	p := t.problem
	value[0] = p.time.Recompute(p.schedule, p.transform(node)).Span

	atomic.AddUint32(&t.ec, 1)
}

func (t *delayTarget) InputsOutputs() (uint32, uint32) {
	return t.problem.zc, 1
}

func (t *delayTarget) Evaluations() uint32 {
	return t.ec
}

func (t *delayTarget) String() string {
	ic, oc := t.InputsOutputs()
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", ic, oc)
}
