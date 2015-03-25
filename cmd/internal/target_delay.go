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
	return t.problem.nz, 2
}

func (t *delayTarget) Compute(node []float64, value []float64) {
	p := t.problem
	value[0] = p.time.Recompute(p.schedule, p.transform(node)).Span
	value[1] = value[0] * value[0]
}

func (t *delayTarget) Refine(node, surplus []float64, volume float64) float64 {
	return Refine(t, node, surplus, volume)
}

func (t *delayTarget) Monitor(level, np, na uint) {
	Monitor(t, level, np, na)
}

func (t *delayTarget) Generate(ns uint) []float64 {
	return Generate(t, ns)
}
