package internal

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

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
