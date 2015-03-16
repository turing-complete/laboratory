package internal

import (
	"math"
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestCombine(t *testing.T) {
	m, n := uint(4), uint(3)

	A := []float64{
		+0, +1, -2, +0,
		-1, -2, +0, +1,
		+1, +1, +0, +2,
	}

	test := func(x, y []float64) {
		z := make([]float64, m)
		combine(A, x, z, m, n)
		assert.Equal(z, y, t)
	}

	inf := math.Inf(1)

	test([]float64{1, 2, 1}, []float64{-1, -2, -2, 4})
	test([]float64{inf, 2, 1}, []float64{-1, inf, -inf, 4})
	test([]float64{1, -inf, 1}, []float64{inf, inf, -2, -inf})
	test([]float64{1, 2, inf}, []float64{inf, inf, -2, inf})
	test([]float64{inf, 2, -inf}, []float64{-inf, -4, -inf, -inf})
}

func TestLocate(t *testing.T) {
	line := []float64{0, 0.2, 0.4, 0.6, 0.8, 1}

	test := func(l, r float64, i, j uint) {
		goti, gotj := locate(l, r, line)
		assert.Equal(goti, i, t)
		assert.Equal(gotj, j, t)
	}

	test(0.0, 1.0, 0, 6)
	test(0.1, 0.9, 0, 6)
	test(0.1, 0.8, 0, 5)
	test(0.2, 0.9, 1, 6)
	test(0.2, 0.8, 1, 5)
	test(0.3, 0.5, 1, 4)
}

func TestSlice(t *testing.T) {
	data := []float64{
		0, 1, 2, 3, 4, 5,
		6, 7, 8, 9, 8, 7,
		6, 5, 4, 3, 2, 1,
	}

	assert.Equal(slice(data, []uint{1, 3}, 6), []float64{1, 3, 7, 9, 5, 3}, t)
	assert.Equal(slice(data, []uint{0, 5}, 6), []float64{0, 5, 6, 7, 6, 1}, t)
}
