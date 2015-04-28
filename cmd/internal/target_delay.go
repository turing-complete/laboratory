package internal

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
	return String(t)
}

func (t *delayTarget) Config() *TargetConfig {
	return t.config
}

func (t *delayTarget) Dimensions() (uint, uint) {
	return t.problem.model.nz, 2
}

func (t *delayTarget) Compute(node []float64, value []float64) {
	s, m := t.problem.system, t.problem.model

	value[0] = s.computeSchedule(m.transform(node)).Span
	value[1] = value[0] * value[0]
}

func (t *delayTarget) Score(node, surplus []float64, volume float64) float64 {
	return Score(t, node, surplus, volume)
}

func (t *delayTarget) Monitor(level, na, nr, nc uint) {
	Monitor(t, level, na, nr, nc)
}
