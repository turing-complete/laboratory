package internal

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestNewConfig(t *testing.T) {
	config, _ := NewConfig("fixtures/004_040_profile.json")

	assert.Equal(config.System, "fixtures/004_040.tgff", t)
	assert.Equal(config.Temperature.Floorplan, "fixtures/004.flp", t)
	assert.Equal(config.Probability.Marginal, "Beta(1, 1)", t)
	assert.Equal(config.Assessment.Samples, uint(10000), t)
}
