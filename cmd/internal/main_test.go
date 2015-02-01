package internal

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestSetup(t *testing.T) {
	config, _ := loadConfig("fixtures/002_020.json")
	problem, _ := newProblem(config)

	target, _, _ := Setup(problem)
	temperatureTarget := target.(*temperatureTarget)

	assert.Equal(temperatureTarget.sc, uint32(29100), t)
}
