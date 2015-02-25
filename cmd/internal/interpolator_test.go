package internal

import (
	"testing"

	"github.com/ready-steady/numeric/grid/newcot"
	"github.com/ready-steady/support/assert"
)

func TestInterpolatorCompute(t *testing.T) {
	config, _ := NewConfig("fixtures/002_020_profile.json")
	problem, _ := NewProblem(config)
	target, _ := NewTarget(problem)
	interpolator, _ := NewInterpolator(problem, target)
	surrogate := interpolator.Compute(target.Evaluate)

	ic, oc := target.Inputs(), target.Outputs()
	nc := surrogate.Nodes

	assert.Equal(nc, uint(4127), t)

	grid := newcot.NewOpen(ic)
	nodes := grid.ComputeNodes(surrogate.Indices)

	values := make([]float64, nc*oc)
	for i := uint(0); i < nc; i++ {
		target.Evaluate(nodes[i*ic:(i+1)*ic], values[i*oc:(i+1)*oc], nil)
	}

	assert.EqualWithin(values, interpolator.Evaluate(surrogate, nodes), 1e-15, t)
}

func BenchmarkInterpolatorCompute(b *testing.B) {
	config, _ := NewConfig("fixtures/002_020_slice.json")
	problem, _ := NewProblem(config)

	for i := 0; i < b.N; i++ {
		target, _ := NewTarget(problem)
		interpolator, _ := NewInterpolator(problem, target)
		interpolator.Compute(target.Evaluate)
	}
}
