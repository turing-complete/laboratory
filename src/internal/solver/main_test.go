package solver

import (
	"testing"

	"github.com/ready-steady/adapt/grid/newcot"
	"github.com/ready-steady/assert"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/target"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

func TestSolverCompute(t *testing.T) {
	config, _ := config.New("fixtures/002_020_temperature.json")
	system, _ := system.New(&config.System)
	uncertainty, _ := uncertainty.New(system, &config.Uncertainty)

	target, _ := target.New(system, uncertainty, &config.Target)
	ni, no := target.Dimensions()

	solver, _ := New(ni, no, &config.Solver)
	solution := solver.Compute(target)

	nc := solution.Surrogate.Nodes

	assert.Equal(nc, uint(857), t)

	grid := newcot.NewOpen(ni)
	nodes := grid.Compute(solution.Surrogate.Indices)

	values := make([]float64, nc*no)
	for i := uint(0); i < nc; i++ {
		target.Compute(nodes[i*ni:(i+1)*ni], values[i*no:(i+1)*no])
	}

	assert.EqualWithin(values, solver.Evaluate(solution, nodes), 1e-15, t)
}
