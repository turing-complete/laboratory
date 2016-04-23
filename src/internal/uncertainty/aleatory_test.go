package uncertainty

import (
	"math"
	"testing"

	"github.com/ready-steady/assert"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

func TestNewAleatory(t *testing.T) {
	const (
		nt = 10
		σ  = 0.2
	)

	config, _ := config.New("fixtures/001_010.json")
	system, _ := system.New(&config.System)
	reference := system.ReferenceTime()
	uncertainty, _ := newAleatory(system, reference, &config.Uncertainty.Time)

	for i := 0; i < nt; i++ {
		min, max := (1.0-σ)*reference[i], (1.0+σ)*reference[i]
		assert.EqualWithin(uncertainty.marginals[i].Decumulate(0.0), min, 1e-15, t)
		assert.EqualWithin(uncertainty.marginals[i].Decumulate(1.0), max, 1e-15, t)
	}
}

func TestMultiply(t *testing.T) {
	m, n := uint(4), uint(3)

	A := []float64{
		+0.0, +1.0, -2.0, +0.0,
		-1.0, -2.0, +0.0, +1.0,
		+1.0, +1.0, +0.0, +2.0,
	}

	test := func(x, y []float64) {
		z := make([]float64, m)
		multiply(A, x, z, m, n)
		assert.Equal(z, y, t)
	}

	inf := math.Inf(1.0)

	test([]float64{1.0, 2.0, 1.0}, []float64{-1.0, -2.0, -2.0, 4.0})
	test([]float64{inf, 2.0, 1.0}, []float64{-1.0, inf, -inf, 4.0})
	test([]float64{1.0, -inf, 1.0}, []float64{inf, inf, -2.0, -inf})
	test([]float64{1.0, 2.0, inf}, []float64{inf, inf, -2.0, inf})
	test([]float64{inf, 2.0, -inf}, []float64{-inf, -4.0, -inf, -inf})
}
