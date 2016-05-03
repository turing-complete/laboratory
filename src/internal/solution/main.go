package solution

import (
	"errors"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/quantity"

	interpolation "github.com/ready-steady/adapt/algorithm/external"
	algorithm "github.com/ready-steady/adapt/algorithm/hybrid"
	basis "github.com/ready-steady/adapt/basis/polynomial"
	grid "github.com/ready-steady/adapt/grid/equidistant"
)

type Solution struct {
	algorithm.Algorithm

	config *config.Solution
	guide  algorithm.Guide
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

	var agrid interface {
		algorithm.Grid
		algorithm.Guide
	}
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
		Algorithm: *algorithm.New(ni, no, agrid, abasis),

		config: config,
		guide:  agrid,
	}, nil
}

func (self *Solution) Compute(target, reference quantity.Quantity) *Surrogate {
	strategy := newStrategy(target, reference, self.guide, self.config)
	surrogate := self.Algorithm.Compute(target.Compute, strategy)
	return &Surrogate{
		Surrogate:  *surrogate,
		Statistics: Statistics{strategy.active},
	}
}

func (self *Solution) Evaluate(surrogate *Surrogate, nodes []float64) []float64 {
	return self.Algorithm.Evaluate(&surrogate.Surrogate, nodes)
}
