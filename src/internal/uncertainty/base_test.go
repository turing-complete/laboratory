package uncertainty

import (
	"math"
	"testing"

	"github.com/ready-steady/assert"
	"github.com/ready-steady/probability/distribution"
)

func TestBaseForwardInverse(t *testing.T) {
	uncertainty := &base{
		tasks: []uint{0, 1, 2},
		lower: []float64{42.0, 42.0, 42.0},
		upper: []float64{42.0, 42.0, 42.0},

		nt: 3,
		nu: 3,
		nz: 2,

		correlation: &correlation{
			C: []float64{
				1.0, 2.0, 3.0,
				4.0, 5.0, 6.0,
			},
			D: []float64{
				6.0, 5.0,
				4.0, 3.0,
				2.0, 1.0,
			},
		},
		marginals: []distribution.Continuous{
			distribution.NewUniform(10.0, 20.0),
			distribution.NewUniform(20.0, 30.0),
			distribution.NewUniform(30.0, 40.0),
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
