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
	return TargetExt{t}.String()
}

func (t *delayTarget) Dimensions() (uint, uint) {
	return t.problem.nz, 2
}

func (t *delayTarget) Compute(node []float64, value []float64) {
	p := t.problem
	value[0] = p.time.Recompute(p.schedule, p.transform(node)).Span
	value[1] = value[0] * value[0]
}

func (t *delayTarget) Refine(surplus []float64) bool {
	return surplus[0] > t.config.Tolerance || -surplus[0] > t.config.Tolerance
}

func (t *delayTarget) Monitor(level, np, na uint) {
	if t.config.Verbose {
		TargetExt{t}.Monitor(level, np, na)
	}
}
