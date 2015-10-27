package config

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestNew(t *testing.T) {
	config, _ := New("fixtures/004_040_profile.json")

	assert.Equal(config.System.Floorplan, "fixtures/004.flp", t)
	assert.Equal(config.System.Configuration, "fixtures/hotspot.config", t)
	assert.Equal(config.System.Specification, "fixtures/004_040.tgff", t)
}
