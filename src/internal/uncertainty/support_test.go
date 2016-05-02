package uncertainty

import (
	"testing"

	"github.com/ready-steady/assert"
	"github.com/ready-steady/linear/decomposition"
	"github.com/ready-steady/linear/matrix"
)

func TestInverse(t *testing.T) {
	m := uint(3)

	A := []float64{
		1.0, 2.0, 3.0,
		2.0, 4.0, 5.0,
		3.0, 5.0, 6.0,
	}
	U := make([]float64, m*m)
	Λ := make([]float64, m)

	err := decomposition.SymmetricEigen(A, U, Λ, m)
	assert.Equal(err, nil, t)

	err = matrix.Invert(A, m)
	assert.Equal(err, nil, t)

	I, err := invert(U, Λ, m)
	assert.Equal(err, nil, t)
	assert.EqualWithin(A, I, 1e-14, t)
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

	test(
		[]float64{1.0, 2.0, 1.0},
		[]float64{-1.0, -2.0, -2.0, 4.0},
	)
	test(
		[]float64{infinity, 2.0, 1.0},
		[]float64{-1.0, infinity, -infinity, 4.0},
	)
	test(
		[]float64{1.0, -infinity, 1.0},
		[]float64{infinity, infinity, -2.0, -infinity},
	)
	test(
		[]float64{1.0, 2.0, infinity},
		[]float64{infinity, infinity, -2.0, infinity},
	)
	test(
		[]float64{infinity, 2.0, -infinity},
		[]float64{-infinity, -4.0, -infinity, -infinity},
	)
}

func TestQuadratic(t *testing.T) {
	m := uint(3)

	A := []float64{
		+0.0, +1.0, -2.0,
		-1.0, -2.0, +0.0,
		+1.0, +1.0, +0.0,
	}

	test := func(x []float64, y float64) {
		assert.Equal(quadratic(A, x, m), y, t)
	}

	test([]float64{1.0, 2.0, 3.0}, -5.0)
	test([]float64{infinity, 2.0, 3.0}, -infinity)
	test([]float64{1.0, infinity, 3.0}, -infinity)
	test([]float64{1.0, 2.0, infinity}, infinity)
	test([]float64{1.0, 1.0, infinity}, -2.0)
}
