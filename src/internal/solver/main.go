package solver

import (
	"errors"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/target"

	interpolation "github.com/ready-steady/adapt/algorithm/external"
	algorithm "github.com/ready-steady/adapt/algorithm/local"
	basis "github.com/ready-steady/adapt/basis/polynomial"
	grid "github.com/ready-steady/adapt/grid/equidistant"
)

type Solver struct {
	algorithm.Interpolator

	strategy func() *strategy
}

type Statistics struct {
	Active []uint
}

type Solution struct {
	interpolation.Surrogate
	Statistics
}

func New(ni, no uint, config *config.Solver) (*Solver, error) {
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

	return &Solver{
		Interpolator: *algorithm.New(ni, no, agrid, abasis),
		strategy:     newStrategy(ni, no, config, agrid),
	}, nil
}

func (self *Solver) Compute(target target.Target) *Solution {
	strategy := self.strategy()
	surrogate := self.Interpolator.Compute(target.Compute, strategy)
	return &Solution{
		Surrogate:  *surrogate,
		Statistics: Statistics{strategy.active},
	}
}

func (self *Solver) Evaluate(solution *Solution, nodes []float64) []float64 {
	return self.Interpolator.Evaluate(&solution.Surrogate, nodes)
}
