package target

import (
	"github.com/ready-steady/adapt"
	"github.com/simulated-reality/laboratory/src/internal/config"
	"github.com/simulated-reality/laboratory/src/internal/problem"
	"github.com/simulated-reality/laboratory/src/internal/support"
)

type profile struct {
	problem *problem.Problem
	config  *config.Target

	coreIndex []uint
	timeIndex []float64
}

func newProfile(p *problem.Problem, c *config.Target) (*profile, error) {
	// The cores of interest.
	coreIndex, err := support.ParseNaturalIndex(c.CoreIndex, 0,
		uint(p.System.Platform.Len())-1)
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
		timeIndex[i] *= p.System.Span()
	}

	target := &profile{
		problem: p,
		config:  c,

		coreIndex: coreIndex,
		timeIndex: timeIndex,
	}

	return target, nil
}

func (t *profile) Dimensions() (uint, uint) {
	nci, nsi := uint(len(t.coreIndex)), uint(len(t.timeIndex))
	return uint(t.problem.Uncertainty.Len()), nsi * nci * 2
}

func (t *profile) Compute(node, value []float64) {
	const (
		ε = 1e-10
	)

	s, u := t.problem.System, t.problem.Uncertainty

	schedule := s.ComputeSchedule(u.Transform(node))
	P, ΔT, timeIndex := s.PartitionPower(schedule, t.timeIndex, ε)
	for i := range timeIndex {
		if timeIndex[i] == 0 {
			panic("the timeline of interest should not contain time 0")
		}
		timeIndex[i]--
	}

	Q := s.ComputeTemperature(P, ΔT)

	coreIndex := t.coreIndex
	nc := uint(s.Platform.Len())
	nci, nsi := uint(len(coreIndex)), uint(len(timeIndex))

	for i, k := uint(0), uint(0); i < nsi; i++ {
		for j := uint(0); j < nci; j++ {
			value[k] = Q[timeIndex[i]*nc+coreIndex[j]]
			value[k+1] = value[k] * value[k]
			k += 2
		}
	}
}

func (t *profile) Monitor(progress *adapt.Progress) {
	if t.config.Verbose {
		Monitor(t, progress)
	}
}

func (t *profile) Score(location *adapt.Location, progress *adapt.Progress) float64 {
	return Score(t, t.config, location, progress)
}

func (t *profile) String() string {
	return String(t)
}
