package internal

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestNewProblem(t *testing.T) {
	config, _ := NewConfig("fixtures/002_020_slice.json")
	problem, _ := NewProblem(config)

	assert.Equal(problem.cc, uint(2), t)
	assert.Equal(problem.tc, uint(20), t)
	assert.Equal(problem.uc, uint(20), t)
	assert.Equal(problem.zc, uint(3), t)

	delay := make([]float64, 20)
	for i := 0; i < 20; i++ {
		assert.Equal(problem.marginals[i].InvCDF(0), 0.0, t)
		delay[i] = problem.marginals[i].InvCDF(1)
	}
	assert.AlmostEqual(delay, []float64{
		0.0020, 0.0006, 0.0076, 0.0062, 0.0004, 0.0038, 0.0006, 0.0062, 0.0036, 0.0056,
		0.0038, 0.0010, 0.0068, 0.0070, 0.0078, 0.0044, 0.0004, 0.0058, 0.0056, 0.0040,
	}, t)

	assert.Equal(len(problem.multiplier), 3*20, t)

	assert.Equal(problem.schedule.Mapping, []uint{
		0, 1, 0, 0, 1, 1, 1, 0, 0, 1,
		1, 0, 0, 0, 0, 1, 1, 1, 1, 1,
	}, t)
	assert.Equal(problem.schedule.Order, []uint{
		0, 1, 2, 9, 12, 16, 18, 14, 17, 13,
		15, 3, 5, 11, 19, 8, 7, 6, 4, 10,
	}, t)
	assert.AlmostEqual(problem.schedule.Start, []float64{
		0.000, 0.010, 0.013, 0.187, 0.265, 0.218, 0.262, 0.260, 0.242, 0.051,
		0.267, 0.237, 0.079, 0.152, 0.113, 0.170, 0.079, 0.141, 0.113, 0.242,
	}, t)
	assert.AlmostEqual(problem.schedule.Span, 0.291, t)
}
