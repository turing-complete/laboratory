package internal

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestNewTarget(t *testing.T) {
	problem, _ := NewProblem("fixtures/002_020.json")

	target, _ := NewTarget(problem)
	temperatureTarget := target.(*temperatureTarget)

	assert.Equal(temperatureTarget.sc, uint32(29100), t)
}
