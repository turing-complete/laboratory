package internal

import (
	"testing"

	"github.com/ready-steady/assert"
	"github.com/simulated-reality/laboratory/internal/config"
)

func TestNewProblem(t *testing.T) {
	config, _ := config.New("fixtures/002_020_profile.json")
	problem, _ := NewProblem(config)

	model := problem.model
	assert.Equal(model.nu, uint(20), t)
	assert.Equal(model.nz, uint(3), t)
	assert.Equal(len(model.correlator), 3*20, t)
}
