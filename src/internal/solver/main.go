package solver

import (
	"errors"
	"fmt"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/target"

	interpolation "github.com/ready-steady/adapt/algorithm/local"
	basis "github.com/ready-steady/adapt/basis/polynomial"
	grid "github.com/ready-steady/adapt/grid/equidistant"
)

type Solver struct {
	interpolation.Interpolator
}

type Solution struct {
	interpolation.Surrogate
}

func New(ni, _ uint, config *config.Solver) (*Solver, error) {
	power := config.Power
	if power == 0 {
		return nil, errors.New("the interpolation power should be positive")
	}
	switch config.Rule {
	case "closed":
		return &Solver{*interpolation.New(grid.NewClosed(ni), basis.NewClosed(ni, power),
			(*interpolation.Config)(&config.Config))}, nil
	case "open":
		return &Solver{*interpolation.New(grid.NewOpen(ni), basis.NewOpen(ni, power),
			(*interpolation.Config)(&config.Config))}, nil
	default:
		return nil, errors.New("the interpolation rule is unknown")
	}
}

func (self *Solver) Compute(target target.Target) *Solution {
	return &Solution{*self.Interpolator.Compute(target, target)}
}

func (self *Solver) Evaluate(solution *Solution, nodes []float64) []float64 {
	return self.Interpolator.Evaluate(&solution.Surrogate, nodes)
}

func (self *Solver) Integrate(solution *Solution) []float64 {
	return self.Interpolator.Integrate(&solution.Surrogate)
}

func (self *Solution) String() string {
	return fmt.Sprintf(`{"inputs": %d, "outputs": %d, "level": %d, "nodes": %d}`,
		self.Inputs, self.Outputs, self.Level, self.Nodes)
}
