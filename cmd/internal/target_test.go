package internal

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestNewTempTarget(t *testing.T) {
	config, _ := loadConfig("fixtures/002_020.json")
	problem, _ := newProblem(config)

	target, err := newTarget(problem)
	assert.Success(err, t)

	tempTarget := target.(*tempTarget)

	assert.Equal(tempTarget.ic, uint32(3+1), t)
	assert.Equal(tempTarget.oc, uint32(1), t)
	assert.Equal(tempTarget.sc, uint32(29100), t)
}
