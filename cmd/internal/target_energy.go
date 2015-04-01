package internal

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
	return String(t)
}

func (t *energyTarget) Config() *TargetConfig {
	return t.config
}

func (t *energyTarget) Dimensions() (uint, uint) {
	return t.problem.nz, 2
}

func (t *energyTarget) Compute(node, value []float64) {
	p := t.problem
	s := p.system

	schedule := s.computeSchedule(p.transform(node))
	time, power := s.computeTime(schedule), s.computePower(schedule)

	value[0] = 0
	for i := range time {
		value[0] += time[i] * power[i]
	}

	value[1] = value[0] * value[0]
}

func (t *energyTarget) Refine(node, surplus []float64, volume float64) float64 {
	return Refine(t, node, surplus, volume)
}

func (t *energyTarget) Monitor(level, np, na uint) {
	Monitor(t, level, np, na)
}
