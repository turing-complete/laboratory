package uncertainty

import (
	"math"
	"testing"

	"github.com/ready-steady/assert"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

func TestNewAleatory001(t *testing.T) {
	config, _ := config.New("fixtures/001_010.json")
	system, _ := system.New(&config.System)
	reference := system.ReferenceTime()
	uncertainty, _ := newAleatory(system, reference, &config.Uncertainty.Time)

	for i := 0; i < 10; i++ {
		min, max := 0.8*reference[i], 1.2*reference[i]
		assert.EqualWithin(uncertainty.marginals[i].InvCDF(0.0), min, 1e-15, t)
		assert.EqualWithin(uncertainty.marginals[i].InvCDF(1.0), max, 1e-15, t)
	}
}

func TestNewAleatory002(t *testing.T) {
	config, _ := config.New("fixtures/002_020.json")
	system, _ := system.New(&config.System)
	reference := system.ReferenceTime()
	uncertainty, _ := newAleatory(system, reference, &config.Uncertainty.Time)

	for i := 0; i < 20; i++ {
		min, max := 0.8*reference[i], 1.2*reference[i]
		assert.EqualWithin(uncertainty.marginals[i].InvCDF(0.0), min, 1e-15, t)
		assert.EqualWithin(uncertainty.marginals[i].InvCDF(1.0), max, 1e-15, t)
	}
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
