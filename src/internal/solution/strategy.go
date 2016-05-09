package solution

import (
	"log"
	"math"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/quantity"

	interpolation "github.com/ready-steady/adapt/algorithm/external"
	algorithm "github.com/ready-steady/adapt/algorithm/hybrid"
)

type strategy struct {
	algorithm.Strategy

	nmax uint

	ns uint
	nn uint

	active []uint
}

func newStrategy(target, reference quantity.Quantity, guide algorithm.Guide,
	config *config.Solution) *strategy {

	ni, no := target.Dimensions()
	return &strategy{
		Strategy: *algorithm.NewStrategy(ni, no, guide, config.MinLevel,
			config.MaxLevel, config.LocalError, config.TotalError),

		target:    target,
		reference: reference,

		nmax: config.MaxEvaluations,
	}
}

func (self *strategy) Done(state *interpolation.State, surrogate *interpolation.Surrogate) bool {
	if self.ns == 0 {
		log.Printf("%5s %15s %15s %15s\n", "Step", "Old Nodes", "New Nodes", "New Level")
	}

	if self.Strategy.Done(state, surrogate) {
		return true
	}

	ni := surrogate.Inputs
	nn := uint(len(state.Indices)) / ni
	if self.nn+nn > self.nmax {
		return true
	}

	level := maxLevel(state.Lindices, ni)

	log.Printf("%5d %15d %15d %15d\n", self.ns, self.nn, nn, level)

	self.ns += 1
	self.nn += nn
	self.active = append(self.active, nn)

	return false
}

func (self *strategy) Score(element *interpolation.Element) float64 {
	return maxAbsolute(element.Surplus) * element.Volume
}

func maxAbsolute(data []float64) (value float64) {
	for i, n := uint(0), uint(len(data)); i < n; i++ {
		value = math.Max(value, math.Abs(data[i]))
	}
	return
}

func maxLevel(lindices []uint64, ni uint) (level uint64) {
	nn := uint(len(lindices)) / ni
	for i := uint(0); i < nn; i++ {
		l := uint64(0)
		for j := uint(0); j < ni; j++ {
			l += lindices[i*ni+j]
		}
		if l > level {
			level = l
		}
	}
	return
}
