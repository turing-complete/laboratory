package solution

import (
	"testing"

	"github.com/ready-steady/assert"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/target"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"

	grid "github.com/ready-steady/adapt/grid/equidistant"
)

func TestSolutionCompute(t *testing.T) {
	config, _ := config.New("fixtures/002_020.json")
	system, _ := system.New(&config.System)
	uncertainty, _ := uncertainty.NewEpistemic(system, &config.Uncertainty)

	target, _ := target.New(system, uncertainty, &config.Target)
	ni, no := target.Dimensions()

	solution, _ := New(ni, no, &config.Solution)
	surrogate := solution.Compute(target)

	nc := surrogate.Surrogate.Nodes

	assert.Equal(nc, uint(490), t)

	grid := grid.NewClosed(ni)
	nodes := grid.Compute(surrogate.Surrogate.Indices)

	values := make([]float64, nc*no)
	for i := uint(0); i < nc; i++ {
		target.Compute(nodes[i*ni:(i+1)*ni], values[i*no:(i+1)*no])
	}

	assert.EqualWithin(values, solution.Evaluate(surrogate, nodes), 1e-15, t)
}
