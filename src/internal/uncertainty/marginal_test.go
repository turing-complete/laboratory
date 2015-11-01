package uncertainty

import (
	"testing"

	"github.com/ready-steady/assert"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

func TestNewMarginal001(t *testing.T) {
	config, _ := config.New("fixtures/001_010_delay.json")
	system, _ := system.New(&config.System)
	uncertainty, _ := NewMarginal(system, &config.Uncertainty)

	delay := make([]float64, 10)
	for i := 0; i < 10; i++ {
		assert.Equal(uncertainty.marginals[i].InvCDF(0), 0.0, t)
		delay[i] = uncertainty.marginals[i].InvCDF(1)
	}
	assert.EqualWithin(delay, []float64{
		0.0058, 0.0032, 0.0050, 0.0058, 0.0074, 0.0060, 0.0032, 0.0048, 0.0030, 0.0060,
	}, 1e-15, t)
}

func TestNewMarginal002(t *testing.T) {
	config, _ := config.New("fixtures/002_020_profile.json")
	system, _ := system.New(&config.System)
	uncertainty, _ := NewMarginal(system, &config.Uncertainty)

	delay := make([]float64, 20)
	for i := 0; i < 20; i++ {
		assert.Equal(uncertainty.marginals[i].InvCDF(0), 0.0, t)
		delay[i] = uncertainty.marginals[i].InvCDF(1)
	}
	assert.EqualWithin(delay, []float64{
		0.0020, 0.0006, 0.0076, 0.0062, 0.0004, 0.0038, 0.0006, 0.0062, 0.0036, 0.0056,
		0.0038, 0.0010, 0.0068, 0.0070, 0.0078, 0.0044, 0.0004, 0.0058, 0.0056, 0.0040,
	}, 1e-15, t)
}
