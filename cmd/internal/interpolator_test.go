package internal

import (
	"testing"

	"github.com/ready-steady/numeric/grid/newcot"
	"github.com/ready-steady/support/assert"
)

func TestInterpolatorCompute(t *testing.T) {
	config, _ := NewConfig("fixtures/002_020_temperature.json")
	problem, _ := NewProblem(config)
	target, _ := NewTarget(problem)
	interpolator, _ := NewInterpolator(problem, target)
	surrogate := interpolator.Compute(target.Evaluate)

	ni, no := target.Inputs(), target.Outputs()
	nc := surrogate.Nodes

	assert.Equal(nc, uint(111), t)

	grid := newcot.NewOpen(ni)
	nodes := grid.ComputeNodes(surrogate.Indices)

	values := make([]float64, nc*no)
	for i := uint(0); i < nc; i++ {
		target.Evaluate(nodes[i*ni:(i+1)*ni], values[i*no:(i+1)*no], nil)
	}

	assert.EqualWithin(values, interpolator.Evaluate(surrogate, nodes), 1e-15, t)
}
