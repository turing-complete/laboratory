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
	switch config.Rule {
	case "closed":
		grid := grid.NewClosed(ni)
		basis := basis.NewClosed(ni, power)
		strategy := newStrategy(ni, no, config, grid)
		return &Solver{*algorithm.New(ni, no, grid, basis, strategy)}, nil
	case "open":
		grid := grid.NewOpen(ni)
		basis := basis.NewOpen(ni, power)
		strategy := newStrategy(ni, no, config, grid)
		return &Solver{*algorithm.New(ni, no, grid, basis, strategy)}, nil
	default:
		return nil, errors.New("the interpolation rule is unknown")
	}
}

func (self *Solver) Compute(target target.Target) *Solution {
	ni, _ := target.Dimensions()
	active := ([]uint)(nil)
	surrogate := self.Interpolator.Compute(func(nodes, values []float64) {
		active = append(active, uint(len(nodes))/ni)
		target.Compute(nodes, values)
	})
	return &Solution{
		Surrogate:  *surrogate,
		Statistics: Statistics{active},
	}
}

func (self *Solver) Evaluate(solution *Solution, nodes []float64) []float64 {
	return self.Interpolator.Evaluate(&solution.Surrogate, nodes)
}
