package uncertainty

import (
	"testing"

	"github.com/ready-steady/assert"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

func TestNewModal(t *testing.T) {
	config, _ := config.New("fixtures/002_020_profile.json")
	system, _ := system.New(&config.System)
	uncertainty, _ := NewModal(system, &config.Uncertainty)

	assert.Equal(uncertainty.nu, uint(20), t)
	assert.Equal(uncertainty.nz, uint(3), t)
	assert.Equal(len(uncertainty.correlator), 3*20, t)
}
