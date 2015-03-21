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
	return CommonTarget{t}.String()
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
	Δ := CommonTarget{t}.Refine(node, surplus, volume)
	if Δ <= t.config.Tolerance {
		Δ = 0
	}
	return Δ
}

func (t *energyTarget) Monitor(level, np, na uint) {
	if t.config.Verbose {
		CommonTarget{t}.Monitor(level, np, na)
	}
}

func (t *energyTarget) Generate(ns uint) []float64 {
	return CommonTarget{t}.Generate(ns)
}
