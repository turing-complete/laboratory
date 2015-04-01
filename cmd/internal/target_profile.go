package internal

type profileTarget struct {
	problem *Problem
	config  *TargetConfig

	coreIndex []uint
	timeline  []float64
	shift     uint
}

func newProfileTarget(p *Problem, c *TargetConfig) (*profileTarget, error) {
	// The cores of interest.
	coreIndex, err := enumerate(p.system.nc, c.CoreIndex)
	if err != nil {
		return nil, err
	}

	// The time moments of interest.
	timeline, err := subdivide(p.system.schedule.Span, c.TimeStep, c.TimeFraction)
	if err != nil {
		return nil, err
	}

	shift := uint(0)

	// Force the first time moment to be zero.
	if timeline[0] != 0 {
		shift++
		timeline = append([]float64{0}, timeline...)
	}

	// Make sure to have at least three time moments.
	if len(timeline) == 2 {
		shift++
		timeline = []float64{0, timeline[1] / 2, timeline[1]}
	}

	target := &profileTarget{
		problem: p,
		config:  c,

		coreIndex: coreIndex,
		timeline:  timeline,
		shift:     shift,
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
	nci, ns := uint(len(t.coreIndex)), uint(len(t.timeline))-t.shift
	return t.problem.nz, ns * nci * 2
}

func (t *profileTarget) Compute(node, value []float64) {
	p := t.problem
	s := p.system

	schedule := s.computeSchedule(p.transform(node))
	Q, _, err := s.temperature.Compute(s.power.Process(schedule), t.timeline)
	if err != nil {
		panic("cannot compute a temperature profile")
	}

	coreIndex := t.coreIndex
	nc, nci, ns := s.nc, uint(len(coreIndex)), uint(len(t.timeline))-t.shift

	Q = Q[t.shift*nc:]

	for i, k := uint(0), uint(0); i < ns; i++ {
		for j := uint(0); j < nci; j++ {
			value[k] = Q[i*nc+coreIndex[j]]
			value[k+1] = value[k] * value[k]
			k += 2
		}
	}
}

func (t *profileTarget) Refine(node, surplus []float64, volume float64) float64 {
	return Refine(t, node, surplus, volume)
}

func (t *profileTarget) Monitor(level, np, na uint) {
	Monitor(t, level, np, na)
}
