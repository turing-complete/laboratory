package internal

type energyTarget struct {
	problem *Problem
}

func newEnergyTarget(p *Problem) Target {
	return &energyTarget{p}
}

func (t *energyTarget) InputsOutputs() (uint32, uint32) {
	return t.problem.zc, t.problem.cc
}

func (t *energyTarget) Evaluate(node, value []float64, _ []uint64) {
	p := t.problem

	cores, tasks := p.platform.Cores, p.application.Tasks
	schedule := p.time.Recompute(p.schedule, p.transform(node))

	value[0] = 0
	for i := range tasks {
		value[0] += (schedule.Finish[i] - schedule.Start[i]) *
			cores[uint32(schedule.Mapping[i])].Power[tasks[i].Type]
	}
}
