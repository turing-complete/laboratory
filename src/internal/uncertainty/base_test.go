package uncertainty

import (
	"math"
	"testing"

	"github.com/ready-steady/assert"
	"github.com/ready-steady/probability"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

func TestBaseForwardInverse(t *testing.T) {
	uncertainty := &base{
		tasks: []uint{0, 1, 2},
		lower: []float64{42.0, 42.0, 42.0},
		upper: []float64{42.0, 42.0, 42.0},

		nt: 3,
		nu: 3,
		nz: 2,

		correlator: []float64{
			1.0, 2.0, 3.0,
			4.0, 5.0, 6.0,
		},
		decorrelator: []float64{
			6.0, 5.0,
			4.0, 3.0,
			2.0, 1.0,
		},
		marginals: []probability.Distribution{
			probability.NewUniform(10.0, 20.0),
			probability.NewUniform(20.0, 30.0),
			probability.NewUniform(30.0, 40.0),
		},
	}

	forward := uncertainty.Forward([]float64{18.0, 21.0, 36.0})
	assert.EqualWithin(forward, []float64{
		6.664804998759882e-01,
		7.313162037785672e-01,
	}, 1e-14, t)

	inverse := uncertainty.Inverse([]float64{0.45, 0.65})
	assert.EqualWithin(inverse, []float64{
		1.921556679782504e+01,
		2.953060310728164e+01,
		3.973501094321997e+01,
	}, 1e-14, t)
}

func TestBasePassThrough(t *testing.T) {
	const (
		nt = 10
		σ  = 0.2
	)

	config, _ := config.New("fixtures/001_010_epistemic.json")
	system, _ := system.New(&config.System)
	reference := system.ReferenceTime()
	uncertainty, _ := newBase(system, reference, &config.Uncertainty.Time)

	point := make([]float64, nt)
	value := make([]float64, nt)
	for i := 0; i < nt; i++ {
		α := float64(i) / (nt - 1)
		point[i] = α
		value[i] = (1.0 - σ + 2.0*σ*α) * reference[i]
	}

	assert.EqualWithin(uncertainty.Inverse(point), value, 1e-15, t)
}

func TestBaseMarginals(t *testing.T) {
	const (
		nt = 10
		σ  = 0.2
	)

	config, _ := config.New("fixtures/001_010_aleatory.json")
	system, _ := system.New(&config.System)
	reference := system.ReferenceTime()
	uncertainty, _ := newBase(system, reference, &config.Uncertainty.Time)

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
