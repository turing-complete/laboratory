package target

import (
	"fmt"
	"log"
	"math"

	"github.com/ready-steady/adapt"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type base struct {
	system      *system.System
	config      *config.Target
	uncertainty uncertainty.Uncertainty

	ni uint
	no uint
}

func (b *base) Dimensions() (uint, uint) {
	return b.ni, b.no
}

func (_ *base) Monitor(progress *adapt.Progress) {
	if progress.Iteration == 0 {
		log.Printf("%5s %10s %15s %15s %15s\n", "Level", "Iteration",
			"Accepted Nodes", "Rejected Nodes", "Current Nodes")
	}
	log.Printf("%5d %10d %15d %15d %15d\n", progress.Level, progress.Iteration,
		progress.Accepted, progress.Rejected, progress.Current)
}

func (b *base) Score(location *adapt.Location, progress *adapt.Progress) float64 {
	config := b.config

	nj := uint(len(config.Importance))

	score, reject, refine := 0.0, true, false
	for i := uint(0); i < b.no; i++ {
		j := i % nj

		if config.Importance[j] == 0 {
			continue
		}

		s := math.Abs(location.Surplus[i])
		if s >= config.Rejection[j] {
			reject = false
		}
		if s > config.Refinement[j] {
			refine = true
		}

		score += config.Importance[j] * s
	}

	if reject {
		return -1
	}
	if !refine {
		return 0
	}

	return score
}

func (b *base) String() string {
	return fmt.Sprintf(`{"inputs": %d, "outputs": %d}`, b.ni, b.no)
}
