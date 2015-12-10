package target

import (
	"fmt"
	"log"
	"math"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"

	interpolation "github.com/ready-steady/adapt/algorithm/local"
)

type base struct {
	system *system.System
	config *config.Target

	ni uint
	no uint
}

func newBase(system *system.System, config *config.Target, ni, no uint) (base, error) {
	return base{system: system, config: config, ni: ni, no: no}, nil
}

func (self *base) Dimensions() (uint, uint) {
	return self.ni, self.no
}

func (_ *base) Monitor(progress *interpolation.Progress) {
	if progress.Level == 0 {
		log.Printf("%5s %15s %15s %15s\n",
			"Level", "Active Nodes", "Passive Nodes", "Refined Nodes")
	}
	log.Printf("%5d %15d %15d %15d\n",
		progress.Level, progress.Active, progress.Passive, progress.Refined)
}

func (self *base) Score(location *interpolation.Location) float64 {
	score := 0.0
	for i := uint(0); i < self.no; i++ {
		score += math.Abs(location.Surplus[i] * location.Volume)
	}
	if score < self.config.Refinement {
		score = 0.0
	}
	return score
}

func (self *base) String() string {
	return fmt.Sprintf(`{"inputs": %d, "outputs": %d}`, self.ni, self.no)
}
