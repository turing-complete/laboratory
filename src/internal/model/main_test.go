package model

import (
	"testing"

	"github.com/ready-steady/assert"
	"github.com/simulated-reality/laboratory/src/internal/config"
	"github.com/simulated-reality/laboratory/src/internal/system"
)

func TestNew(t *testing.T) {
	config, _ := config.New("fixtures/002_020_profile.json")
	system, _ := system.New(&config.System)
	model, _ := New(&config.Probability, system)

	assert.Equal(model.nu, uint(20), t)
	assert.Equal(model.nz, uint(3), t)
	assert.Equal(len(model.correlator), 3*20, t)
}
