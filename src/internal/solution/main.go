package solution

import (
	"errors"

	"github.com/ready-steady/adapt/algorithm"
	"github.com/ready-steady/adapt/algorithm/hybrid"
	"github.com/ready-steady/adapt/basis/polynomial"
	"github.com/ready-steady/adapt/grid"
	"github.com/ready-steady/adapt/grid/equidistant"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/quantity"
)

type Solution struct {
	hybrid.Algorithm

	config *config.Solution
	grid   interface {
		hybrid.Guide
		grid.Parenter
	}
}

type Statistics struct {
	Active []uint
}

type Surrogate struct {
	algorithm.Surrogate
	Statistics
}

func New(ni, no uint, config *config.Solution) (*Solution, error) {
	power := config.Power
	if power == 0 {
		return nil, errors.New("the interpolation power should be positive")
	}

	var agrid interface {
		hybrid.Grid
		hybrid.Guide
		grid.Parenter
	}
	var abasis hybrid.Basis
	switch config.Rule {
	case "closed":
		agrid = equidistant.NewClosed(ni)
		abasis = polynomial.NewClosed(ni, power)
	case "open":
		agrid = equidistant.NewOpen(ni)
		abasis = polynomial.NewOpen(ni, power)
	default:
		return nil, errors.New("the interpolation rule is unknown")
	}

	return &Solution{
		Algorithm: *hybrid.New(ni, no, agrid, abasis),

		config: config,
		grid:   agrid,
	}, nil
}

func (self *Solution) Compute(target, reference quantity.Quantity) *Surrogate {
	strategy := newStrategy(target, reference, self.grid, self.config)
	surrogate := self.Algorithm.Compute(target.Compute, strategy)
	return &Surrogate{
		Surrogate:  *surrogate,
		Statistics: Statistics{strategy.active},
	}
}

func (self *Solution) Evaluate(surrogate *Surrogate, nodes []float64) []float64 {
	return self.Algorithm.Evaluate(&surrogate.Surrogate, nodes)
}

func (self *Solution) Validate(surrogate *Surrogate) bool {
	return algorithm.Validate(surrogate.Indices, surrogate.Inputs, self.grid)
}
