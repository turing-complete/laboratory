package internal

type delayTarget struct {
	problem *Problem
}

func newDelayTarget(p *Problem) Target {
	return &delayTarget{p}
}

func (t *delayTarget) InputsOutputs() (uint32, uint32) {
	return t.problem.zc, t.problem.cc
}

func (t *delayTarget) Evaluate(node []float64, value []float64, _ []uint64) {
	p := t.problem
	value[0] = p.time.Recompute(p.schedule, p.transform(node)).Span
}
