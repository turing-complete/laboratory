package uncertainty

import (
	"math"
	"testing"

	"github.com/ready-steady/assert"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

func TestNewMarginal001(t *testing.T) {
	config, _ := config.New("fixtures/001_010_delay.json")
	system, _ := system.New(&config.System)
	uncertainty, _ := newMarginal(system, system.ReferenceTime(), &config.Target.Uncertainty)

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
	uncertainty, _ := newMarginal(system, system.ReferenceTime(), &config.Target.Uncertainty)

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

func TestMultiply(t *testing.T) {
	m, n := uint(4), uint(3)

	A := []float64{
		+0, +1, -2, +0,
		-1, -2, +0, +1,
		+1, +1, +0, +2,
	}

	test := func(x, y []float64) {
		z := make([]float64, m)
		multiply(A, x, z, m, n)
		assert.Equal(z, y, t)
	}

	inf := math.Inf(1)

	test([]float64{1, 2, 1}, []float64{-1, -2, -2, 4})
	test([]float64{inf, 2, 1}, []float64{-1, inf, -inf, 4})
	test([]float64{1, -inf, 1}, []float64{inf, inf, -2, -inf})
	test([]float64{1, 2, inf}, []float64{inf, inf, -2, inf})
	test([]float64{inf, 2, -inf}, []float64{-inf, -4, -inf, -inf})
}
