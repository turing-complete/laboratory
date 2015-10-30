package solver

import (
	"errors"
	"fmt"

	"github.com/ready-steady/adapt"
	"github.com/ready-steady/adapt/basis/linhat"
	"github.com/ready-steady/adapt/grid/newcot"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/problem"
	"github.com/turing-complete/laboratory/src/internal/target"
)

type Solver struct {
	adapt.Interpolator
}

type Solution struct {
	adapt.Surrogate
}

func New(problem *problem.Problem, target target.Target,
	config *config.Interpolation) (*Solver, error) {

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

func (s *Solver) Compute(target target.Target) *Solution {
	return &Solution{*s.Interpolator.Compute(target)}
}

func (s *Solver) Evaluate(solution *Solution, nodes []float64) []float64 {
	return s.Interpolator.Evaluate(&solution.Surrogate, nodes)
}

func (s *Solver) Integrate(solution *Solution) []float64 {
	return s.Interpolator.Integrate(&solution.Surrogate)
}

func (s *Solution) String() string {
	return fmt.Sprintf(`{"inputs": %d, "outputs": %d, "level": %d, "nodes": %d}`,
		s.Inputs, s.Outputs, s.Level, s.Nodes)
}
