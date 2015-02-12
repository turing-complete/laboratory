package internal

import (
	"fmt"

	"github.com/ready-steady/simulation/power"
	"github.com/ready-steady/simulation/temperature"

	"../../pkg/pool"
)

type profileTarget struct {
	problem *Problem

	sc uint32

	power       *power.Power
	temperature *temperature.Temperature

	pool *pool.Pool
}

type profileData struct {
	P []float64
	S []float64
	Q []float64
}

func newProfileTarget(p *Problem) (Target, error) {
	const (
		poolCapacity = 100
		MaxUInt16    = ^uint16(0)
	)

	c := &p.Config

	power, err := power.New(p.platform, p.application, c.TempAnalysis.TimeStep)
	if err != nil {
		return nil, err
	}

	temperature, err := temperature.New(temperature.Config(c.TempAnalysis))
	if err != nil {
		return nil, err
	}

	cc, sc := p.cc, uint32(p.schedule.Span/c.TempAnalysis.TimeStep)
	nc := temperature.Nodes

	if cc*sc > uint32(MaxUInt16) {
		panic("The number of outputs is too large.")
	}

	target := &profileTarget{
		problem: p,

		sc: sc,

		power:       power,
		temperature: temperature,

		pool: pool.New(poolCapacity, func() interface{} {
			return &profileData{
				P: make([]float64, cc*sc),
				S: make([]float64, nc*sc),
				Q: make([]float64, cc*sc),
			}
		}),
	}

	return target, nil
}

func (t *profileTarget) Inputs() uint32 {
	return t.problem.zc
}

func (t *profileTarget) Outputs() uint32 {
	return t.sc * uint32(len(t.problem.Config.CoreIndex))
}

func (t *profileTarget) Pseudos() uint32 {
	return 0
}

func (t *profileTarget) String() string {
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", t.Inputs(), t.Outputs())
}

func (t *profileTarget) Evaluate(node, value []float64, _ []uint64) {
	p := t.problem

	coreIndex := p.Config.CoreIndex

	cc, occ, sc := p.cc, uint32(len(coreIndex)), t.sc

	data := t.pool.Get().(*profileData)

	Q := data.Q

	t.power.Compute(p.time.Recompute(p.schedule, p.transform(node)), data.P, sc)
	t.temperature.ComputeTransient(data.P, Q, data.S, sc)

	for i := uint32(0); i < sc; i++ {
		for j := uint32(0); j < occ; j++ {
			value[i*cc+j] = Q[i*cc+uint32(coreIndex[j])]
		}
	}

	t.pool.Put(data)
}

func (t *profileTarget) Progress(level uint8, activeNodes, totalNodes uint32) {
	if level == 0 {
		t.problem.Printf("%10s %15s %15s\n",
			"Level", "Passive Nodes", "Active Nodes")
	}

	passiveNodes := totalNodes - activeNodes
	t.problem.Printf("%10d %15d %15d\n", level, passiveNodes, activeNodes)
}
