package internal

import (
	"github.com/ready-steady/adapt"
	"github.com/simulated-reality/laboratory/internal/config"
	"github.com/simulated-reality/laboratory/internal/support"
)

type profileTarget struct {
	problem *Problem
	config  *config.Target

	coreIndex []uint
	timeIndex []float64
}

func newProfileTarget(p *Problem, c *config.Target) (*profileTarget, error) {
	// The cores of interest.
	coreIndex, err := support.ParseNaturalIndex(c.CoreIndex, 0, p.system.nc-1)
	if err != nil {
		return nil, err
	}

	// The time moments of interest.
	timeIndex, err := support.ParseRealIndex(c.TimeIndex, 0, 1)
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

func (t *profileTarget) Monitor(progress *adapt.Progress) {
	if t.config.Verbose {
		Monitor(t, progress)
	}
}

func (t *profileTarget) Score(location *adapt.Location, progress *adapt.Progress) float64 {
	return Score(t, t.config, location, progress)
}
