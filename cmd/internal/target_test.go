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

	heatTarget := target.(*heatTarget)

	assert.Equal(heatTarget.ic, uint32(3+1), t)
	assert.Equal(heatTarget.oc, uint32(2), t)
	assert.Equal(heatTarget.sc, uint32(29100), t)
}
