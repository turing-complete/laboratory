package internal

import (
	"testing"

	"github.com/ready-steady/numeric/grid/newcot"
	"github.com/ready-steady/support/assert"
)

func TestInterpolatorCompute(t *testing.T) {
	problem, _ := NewProblem("fixtures/002_020_profile.json")
	target, _ := NewTarget(problem)
	interpolator, _ := NewInterpolator(problem, target)
	surrogate := interpolator.Compute(target.Evaluate)

	ic, oc := target.InputsOutputs()

	grid := newcot.NewOpen(uint16(ic))
	nodes := grid.ComputeNodes(surrogate.Indices)

	nc := uint32(len(nodes)) / ic

	values := make([]float64, nc*oc)
	for i := uint32(0); i < nc; i++ {
		target.Evaluate(nodes[i*ic:(i+1)*ic], values[i*oc:(i+1)*oc], nil)
	}

	assert.AlmostEqual(values, interpolator.Evaluate(surrogate, nodes), t)
}

func BenchmarkInterpolatorCompute(b *testing.B) {
	problem, _ := NewProblem("fixtures/002_020_slice.json")

	for i := 0; i < b.N; i++ {
		target, _ := NewTarget(problem)
		interpolator, _ := NewInterpolator(problem, target)
		interpolator.Compute(target.Evaluate)
	}
}
