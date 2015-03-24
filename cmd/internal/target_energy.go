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
	return GenericTarget{t}.String()
}

func (t *energyTarget) Config() *TargetConfig {
	return t.config
}

func (t *energyTarget) Dimensions() (uint, uint) {
	return t.problem.nz, 2
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

	value[1] = value[0] * value[0]
}

func (t *energyTarget) Refine(node, surplus []float64, volume float64) float64 {
	return GenericTarget{t}.Refine(node, surplus, volume)
}

func (t *energyTarget) Monitor(level, np, na uint) {
	GenericTarget{t}.Monitor(level, np, na)
}

func (t *energyTarget) Generate(ns uint) []float64 {
	return GenericTarget{t}.Generate(ns)
}
