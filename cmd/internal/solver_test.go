package internal

import (
	"testing"

	"github.com/ready-steady/numeric/grid/newcot"
	"github.com/ready-steady/support/assert"
)

func TestSolverCompute(t *testing.T) {
	config, _ := NewConfig("fixtures/002_020_profile.json")
	problem, _ := NewProblem(config)
	target, _ := NewTarget(problem)
	solver, _ := NewSolver(problem, target)
	solution := solver.Compute(target)
	surrogate := &solution.Surrogate

	ni, no := target.Dimensions()
	nc := surrogate.Nodes

	assert.Equal(nc, uint(89), t)

	grid := newcot.NewOpen(ni)
	nodes := grid.Compute(surrogate.Indices)

	values := make([]float64, nc*no)
	for i := uint(0); i < nc; i++ {
		target.Compute(nodes[i*ni:(i+1)*ni], values[i*no:(i+1)*no])
	}

	assert.EqualWithin(values, solver.Evaluate(solution, nodes), 1e-15, t)
}
