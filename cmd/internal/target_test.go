package internal

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestNewTarget(t *testing.T) {
	config, _ := NewConfig("fixtures/002_020_slice.json")
	problem, _ := NewProblem(config)

	target, _ := NewTarget(problem)
	sliceTarget := target.(*sliceTarget)

	assert.Equal(sliceTarget.interval, []float64{0, 0.291}, t)
}
