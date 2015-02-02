package internal

import (
	"testing"
)

func BenchmarkInterpolatorCompute(b *testing.B) {
	problem, _ := NewProblem("fixtures/002_020.json")

	for i := 0; i < b.N; i++ {
		target, _ := NewTarget(problem)
		interpolator, _ := NewInterpolator(problem, target)
		interpolator.Compute(target.Evaluate)
	}
}
