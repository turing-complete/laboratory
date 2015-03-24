package internal

import (
	"github.com/ready-steady/simulation/power"
	"github.com/ready-steady/simulation/temperature/numeric"
	"github.com/ready-steady/sort"
)

type switchTarget struct {
	problem *Problem
	config  *TargetConfig

	power       *power.Power
	temperature *numeric.Temperature

	cores []uint
}

func newSwitchTarget(p *Problem, tac *TargetConfig,
	tec *TemperatureConfig) (*switchTarget, error) {

	power := power.New(p.platform, p.application)
	temperature, err := numeric.New(&tec.Config)
	if err != nil {
		return nil, err
	}

	// The cores of interest.
	cores, err := enumerate(p.nc, tac.CoreIndex)
	if err != nil {
		return nil, err
	}

	target := &switchTarget{
		problem: p,
		config:  tac,

		power:       power,
		temperature: temperature,

		cores: cores,
	}

	return target, nil
}

func (t *switchTarget) String() string {
	return GenericTarget{t}.String()
}

func (t *switchTarget) Config() *TargetConfig {
	return t.config
}

func (t *switchTarget) Dimensions() (uint, uint) {
	return t.problem.nz, 2 * t.problem.nt * (1 + uint(len(t.cores)))
}

func (t *switchTarget) Compute(node, value []float64) {
	p := t.problem

	cores := t.cores
	nc, nt, nci := p.nc, p.nt, uint(len(cores))

	schedule := p.time.Recompute(p.schedule, p.transform(node))

	timeline := make([]float64, 1+nt) // +1 to start from time 0
	copy(timeline[1:], schedule.Finish)
	_, order := sort.Quick(timeline[1:])

	Q, _, err := t.temperature.Compute(t.power.Process(schedule), timeline)
	if err != nil {
		panic("cannot compute a temperature profile")
	}

	timeline, Q = timeline[1:], Q[nc:] // +1 to exclude time 0

	for i, k := uint(0), uint(0); i < nt; i++ {
		value[k] = schedule.Finish[i]
		value[k+1] = value[k] * value[k]
		k += 2

		for j := uint(0); j < nci; j++ {
			value[k] = Q[order[i]*nc+cores[j]]
			value[k+1] = value[k] * value[k]
			k += 2
		}
	}
}

func (t *switchTarget) Refine(node, surplus []float64, volume float64) float64 {
	return GenericTarget{t}.Refine(node, surplus, volume)
}

func (t *switchTarget) Monitor(level, np, na uint) {
	GenericTarget{t}.Monitor(level, np, na)
}

func (t *switchTarget) Generate(ns uint) []float64 {
	return GenericTarget{t}.Generate(ns)
}
