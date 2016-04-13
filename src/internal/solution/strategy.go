package solution

import (
	"log"

	"github.com/turing-complete/laboratory/src/internal/config"

	interpolation "github.com/ready-steady/adapt/algorithm/external"
	algorithm "github.com/ready-steady/adapt/algorithm/local"
)

type strategy struct {
	algorithm.Strategy

	ns uint
	nn uint

	active []uint
}

func newStrategy(ni, no uint, config *config.Solution, grid algorithm.Grid) func() *strategy {
	return func() *strategy {
		return &strategy{
			Strategy: *algorithm.NewStrategy(ni, no, config.MinLevel,
				config.MaxLevel, config.LocalError, grid),
		}
	}
}

func (self *strategy) Check(state *interpolation.State, surrogate *interpolation.Surrogate) bool {
	if self.ns == 0 {
		log.Printf("%5s %15s %15s\n", "Step", "New Nodes", "Old Nodes")
	}

	nn := uint(len(state.Indices)) / surrogate.Inputs
	log.Printf("%5d %15d %15d\n", self.ns, nn, self.nn)

	self.nn += nn
	self.ns += 1
	self.active = append(self.active, nn)

	return self.Strategy.Check(state, surrogate)
}
