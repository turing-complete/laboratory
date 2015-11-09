package solver

import (
	"errors"
	"fmt"

	"github.com/ready-steady/adapt"
	"github.com/ready-steady/adapt/basis/linhat"
	"github.com/ready-steady/adapt/grid/newcot"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/target"
)

type Solver struct {
	adapt.Interpolator
}

type Solution struct {
	adapt.Surrogate
}

func New(target target.Target, config *config.Solver) (*Solver, error) {
	ni, _ := target.Dimensions()

	var grid adapt.Grid
	var basis adapt.Basis

	switch config.Rule {
	case "closed":
		grid, basis = newcot.NewClosed(ni), linhat.NewClosed(ni)
	case "open":
		grid, basis = newcot.NewOpen(ni), linhat.NewOpen(ni)
	default:
		return nil, errors.New("the interpolation rule is unknown")
	}

	interpolator := adapt.New(grid, basis, (*adapt.Config)(&config.Config))

	return &Solver{*interpolator}, nil
}

func (self *Solver) Compute(target target.Target) *Solution {
	return &Solution{*self.Interpolator.Compute(target)}
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
