package target

import (
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/support"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type profile struct {
	base
	coreIndex []uint
	timeIndex []float64
	time      uncertainty.Uncertainty
	power     uncertainty.Uncertainty
}

func newProfile(system *system.System, config *config.Target) (*profile, error) {
	base, err := newBase(system, config)
	if err != nil {
		return nil, err
	}

	coreIndex, err := support.ParseNaturalIndex(config.CoreIndex, 0, uint(system.Platform.Len())-1)
	if err != nil {
		return nil, err
	}

	timeIndex, err := support.ParseRealIndex(config.TimeIndex, 0, 1)
	if err != nil {
		return nil, err
	}
	if timeIndex[0] == 0 {
		timeIndex = timeIndex[1:]
	}
	for i := range timeIndex {
		timeIndex[i] *= system.Span()
	}

	time, err := uncertainty.New(system, system.ReferenceTime(), &config.Uncertainty)
	if err != nil {
		return nil, err
	}
	power, err := uncertainty.New(system, system.ReferencePower(), &config.Uncertainty)
	if err != nil {
		return nil, err
	}
	base.ni = uint(time.Len() + power.Len())
	base.no = 2 * uint(len(timeIndex)*len(coreIndex))

	return &profile{
		base:      base,
		coreIndex: coreIndex,
		timeIndex: timeIndex,
		time:      time,
		power:     power,
	}, nil
}

func (self *profile) Compute(node, value []float64) {
	const (
		ε = 1e-10
	)

	nt, np := self.time.Len(), self.power.Len()

	time := self.time.Transform(node[:nt])
	power := self.time.Transform(node[nt : nt+np])
	schedule := self.system.ComputeSchedule(time)

	P, ΔT, timeIndex := self.system.PartitionPower(power, schedule, self.timeIndex, ε)
	for i := range timeIndex {
		if timeIndex[i] == 0 {
			panic("the timeline of interest should not contain time 0")
		}
		timeIndex[i]--
	}

	Q := self.system.ComputeTemperature(P, ΔT)

	coreIndex := self.coreIndex
	nc := uint(self.system.Platform.Len())
	nci, nsi := uint(len(coreIndex)), uint(len(timeIndex))

	for i, k := uint(0), uint(0); i < nsi; i++ {
		for j := uint(0); j < nci; j++ {
			value[k] = Q[timeIndex[i]*nc+coreIndex[j]]
			value[k+1] = value[k] * value[k]
			k += 2
		}
	}
}
