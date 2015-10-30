package solver

import (
	"testing"

	"github.com/ready-steady/adapt/grid/newcot"
	"github.com/ready-steady/assert"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/problem"
	"github.com/turing-complete/laboratory/src/internal/target"
)

func TestSolverCompute(t *testing.T) {
	config, _ := config.New("fixtures/002_020_profile.json")
	problem, _ := problem.New(config)
	target, _ := target.New(problem, &config.Target)
	solver, _ := New(problem, target, &config.Solver)
	solution := solver.Compute(target)
	surrogate := &solution.Surrogate

	ni, no := target.Dimensions()
	nc := surrogate.Nodes

	assert.Equal(nc, uint(111), t)

	grid := newcot.NewOpen(ni)
	nodes := grid.Compute(surrogate.Indices)

	values := make([]float64, nc*no)
	for i := uint(0); i < nc; i++ {
		target.Compute(nodes[i*ni:(i+1)*ni], values[i*no:(i+1)*no])
	}

	assert.EqualWithin(values, solver.Evaluate(solution, nodes), 1e-15, t)
}
