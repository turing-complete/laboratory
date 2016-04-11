package target

import (
	"log"
	"math"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"

	interpolation "github.com/ready-steady/adapt/algorithm/external"
)

type base struct {
	system *system.System
	config *config.Target

	ni uint
	no uint

	ns uint
	nn uint
}

func newBase(system *system.System, config *config.Target, ni, no uint) (base, error) {
	return base{system: system, config: config, ni: ni, no: no}, nil
}

func (self *base) Check(state *interpolation.State, _ *interpolation.Surrogate) {
	if self.ns == 0 {
		log.Printf("%5s %15s %15s\n", "", "Done", "More")
	}

	nn := uint(len(state.Indices)) / self.ni

	log.Printf("%5d %15d %15d\n", self.ns, self.nn, nn)

	self.ns += 1
	self.nn += nn
}

func (self *base) Score(element *interpolation.Element) (score float64) {
	for _, value := range element.Surplus {
		score += math.Abs(value)
	}
	score *= element.Volume
	return
}
