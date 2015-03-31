package internal

import (
	"github.com/ready-steady/sort"
)

type switchTarget struct {
	problem *Problem
	config  *TargetConfig

	coreIndex []uint
}

func newSwitchTarget(p *Problem, c *TargetConfig) (*switchTarget, error) {
	// The cores of interest.
	coreIndex, err := enumerate(p.nc, c.CoreIndex)
	if err != nil {
		return nil, err
	}

	target := &switchTarget{
		problem: p,
		config:  c,

		coreIndex: coreIndex,
	}

	return target, nil
}

func (t *switchTarget) String() string {
	return String(t)
}

func (t *switchTarget) Config() *TargetConfig {
	return t.config
}

func (t *switchTarget) Dimensions() (uint, uint) {
	return t.problem.nz, 2 * t.problem.nt * (1 + uint(len(t.coreIndex)))
}

func (t *switchTarget) Compute(node, value []float64) {
	p := t.problem
	s := p.system

	coreIndex := t.coreIndex
	nc, nt, nci := p.nc, p.nt, uint(len(coreIndex))

	schedule := s.computeSchedule(p.transform(node))

	timeline := make([]float64, 1+nt) // +1 to start from time 0
	copy(timeline[1:], schedule.Finish)
	_, order := sort.Quick(timeline[1:])

	Q, _, err := s.temperature.Compute(s.power.Process(schedule), timeline)
	if err != nil {
		panic("cannot compute a temperature profile")
	}

	timeline, Q = timeline[1:], Q[nc:] // +1 to exclude time 0

	for i, k := uint(0), uint(0); i < nt; i++ {
		value[k] = schedule.Finish[i]
		value[k+1] = value[k] * value[k]
		k += 2

		for j := uint(0); j < nci; j++ {
			value[k] = Q[order[i]*nc+coreIndex[j]]
			value[k+1] = value[k] * value[k]
			k += 2
		}
	}
}

func (t *switchTarget) Refine(node, surplus []float64, volume float64) float64 {
	return Refine(t, node, surplus, volume)
}

func (t *switchTarget) Monitor(level, np, na uint) {
	Monitor(t, level, np, na)
}

func (t *switchTarget) Generate(ns uint) []float64 {
	return Generate(t, ns)
}
