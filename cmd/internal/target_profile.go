package internal

type profileTarget struct {
	problem *Problem
	config  *TargetConfig

	coreIndex []uint
	timeIndex []float64
}

func newProfileTarget(p *Problem, c *TargetConfig) (*profileTarget, error) {
	// The cores of interest.
	coreIndex, err := parseNaturalIndex(c.CoreIndex, 0, p.system.nc-1)
	if err != nil {
		return nil, err
	}

	// The time moments of interest.
	timeIndex, err := parseRealIndex(c.TimeIndex, 0, 1)
	if err != nil {
		return nil, err
	}
	if timeIndex[0] == 0 {
		timeIndex = timeIndex[1:]
	}
	for i := range timeIndex {
		timeIndex[i] *= p.system.schedule.Span
	}

	target := &profileTarget{
		problem: p,
		config:  c,

		coreIndex: coreIndex,
		timeIndex: timeIndex,
	}

	return target, nil
}

func (t *profileTarget) String() string {
	return String(t)
}

func (t *profileTarget) Config() *TargetConfig {
	return t.config
}

func (t *profileTarget) Dimensions() (uint, uint) {
	nci, nsi := uint(len(t.coreIndex)), uint(len(t.timeIndex))
	return t.problem.model.nz, nsi * nci * 2
}

func (t *profileTarget) Compute(node, value []float64) {
	const (
		ε = 1e-10
	)

	s, m := t.problem.system, t.problem.model

	schedule := s.computeSchedule(m.transform(node))
	P, ΔT, timeIndex := s.power.Partition(schedule, t.timeIndex, ε)
	for i := range timeIndex {
		if timeIndex[i] == 0 {
			panic("the timeline of interest should not contain time 0")
		}
		timeIndex[i]--
	}

	Q := s.temperature.Compute(P, ΔT)

	coreIndex := t.coreIndex
	nc, nci, nsi := s.nc, uint(len(coreIndex)), uint(len(timeIndex))

	for i, k := uint(0), uint(0); i < nsi; i++ {
		for j := uint(0); j < nci; j++ {
			value[k] = Q[timeIndex[i]*nc+coreIndex[j]]
			value[k+1] = value[k] * value[k]
			k += 2
		}
	}
}

func (t *profileTarget) Score(node, surplus []float64, volume float64) float64 {
	return Score(t, node, surplus, volume)
}

func (t *profileTarget) Monitor(level, na, nr, nc uint) {
	Monitor(t, level, na, nr, nc)
}
