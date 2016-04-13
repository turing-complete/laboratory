package solution

import (
	"errors"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/target"

	interpolation "github.com/ready-steady/adapt/algorithm/external"
	algorithm "github.com/ready-steady/adapt/algorithm/local"
	basis "github.com/ready-steady/adapt/basis/polynomial"
	grid "github.com/ready-steady/adapt/grid/equidistant"
)

type Solution struct {
	algorithm.Interpolator

	strategy func() *strategy
}

type Statistics struct {
	Active []uint
}

type Surrogate struct {
	interpolation.Surrogate
	Statistics
}

func New(ni, no uint, config *config.Solution) (*Solution, error) {
	power := config.Power
	if power == 0 {
		return nil, errors.New("the interpolation power should be positive")
	}

	var agrid algorithm.Grid
	var abasis algorithm.Basis
	switch config.Rule {
	case "closed":
		agrid = grid.NewClosed(ni)
		abasis = basis.NewClosed(ni, power)
	case "open":
		agrid = grid.NewOpen(ni)
		abasis = basis.NewOpen(ni, power)
	default:
		return nil, errors.New("the interpolation rule is unknown")
	}

	return &Solution{
		Interpolator: *algorithm.New(ni, no, agrid, abasis),
		strategy:     newStrategy(ni, no, config, agrid),
	}, nil
}

func (self *Solution) Compute(target target.Target) *Surrogate {
	strategy := self.strategy()
	surrogate := self.Interpolator.Compute(target.Compute, strategy)
	return &Surrogate{
		Surrogate:  *surrogate,
		Statistics: Statistics{strategy.active},
	}
}

func (self *Solution) Evaluate(surrogate *Surrogate, nodes []float64) []float64 {
	return self.Interpolator.Evaluate(&surrogate.Surrogate, nodes)
}
