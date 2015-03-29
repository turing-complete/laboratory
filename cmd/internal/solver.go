package internal

import (
	"errors"
	"fmt"

	"github.com/ready-steady/adhier"
	"github.com/ready-steady/adhier/basis/linhat"
	"github.com/ready-steady/adhier/grid/newcot"
)

type Solver struct {
	adhier.Interpolator
}

type Solution struct {
	adhier.Surrogate
	Expectation []float64
}

func NewSolver(problem *Problem, target Target) (*Solver, error) {
	ni, _ := target.Dimensions()

	var grid adhier.Grid
	var basis adhier.Basis

	switch problem.Config.Interpolation.Rule {
	case "closed":
		grid, basis = newcot.NewClosed(ni), linhat.NewClosed(ni)
	case "open":
		grid, basis = newcot.NewOpen(ni), linhat.NewOpen(ni)
	default:
		return nil, errors.New("the interpolation rule is unknown")
	}

	interpolator := adhier.New(grid, basis,
		(*adhier.Config)(&problem.Config.Interpolation.Config))

	return &Solver{*interpolator}, nil
}

func (s *Solver) Compute(target Target) *Solution {
	surrogate := s.Interpolator.Compute(target)
	target.Monitor(surrogate.Level, surrogate.Nodes, 0)
	return &Solution{
		Surrogate:   *surrogate,
		Expectation: s.Interpolator.Integrate(surrogate),
	}
}

func (s *Solver) Evaluate(solution *Solution, nodes []float64) []float64 {
	return s.Interpolator.Evaluate(&solution.Surrogate, nodes)
}

func (s *Solution) String() string {
	return fmt.Sprintf("Solution{inputs: %d, outputs: %d, level: %d, nodes: %d}",
		s.Inputs, s.Outputs, s.Level, s.Nodes)
}
