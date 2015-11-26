package support

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestParseNaturalIndex(t *testing.T) {
	cases := []struct {
		line   string
		min    uint
		max    uint
		result []uint
	}{
		{"", 0, 10, nil},
		{"[]", 0, 10, []uint{}},
		{"[0, 1, 11]", 0, 10, nil},
		{"[0, 1, 9, 10]", 0, 10, []uint{0, 1, 9, 10}},
		{"[0:10]", 0, 10, []uint{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{"[1:2:10]", 0, 10, []uint{1, 3, 5, 7, 9}},
		{"[0:2:10]", 0, 10, []uint{0, 2, 4, 6, 8, 10}},
		{"[0:5:15]", 0, 10, nil},
		{"[0, 1, end]", 0, 10, []uint{0, 1, 10}},
		{"[0:2:end]", 0, 10, []uint{0, 2, 4, 6, 8, 10}},
		{"[0:end]", 0, 10, []uint{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
	}

	for _, c := range cases {
		result, err := ParseNaturalIndex(c.line, c.min, c.max)
		if c.result != nil {
			assert.Success(err, t)
		}
		assert.Equal(result, c.result, t)
	}
}

func TestParseRealIndex(t *testing.T) {
	cases := []struct {
		line   string
		min    float64
		max    float64
		result []float64
	}{
		{"", 0, 1, nil},
		{"[]", 0, 1, []float64{}},
		{"[0, 0.1, 0.9, 1]", 0, 1, []float64{0.0, 0.1, 0.9, 1.0}},
		{"[0, 0.1, 1.1]", 0, 1, nil},
		{"[0:1]", 0, 1, []float64{0, 1}},
		{"[0.1:0.2:1]", 0, 1, []float64{0.1, 0.3, 0.5, 0.7, 0.9}},
		{"[0:0.2:1]", 0, 1, []float64{0.0, 0.2, 0.4, 0.6, 0.8, 1.0}},
		{"[0:0.5:1.5]", 0, 1, nil},
		{"[0, 0.2, end]", 0, 1, []float64{0.0, 0.2, 1.0}},
		{"[0:0.2:end]", 0, 1, []float64{0.0, 0.2, 0.4, 0.6, 0.8, 1.0}},
		{"[0:end]", 0, 1, []float64{0, 1}},
	}

	for _, c := range cases {
		result, err := ParseRealIndex(c.line, c.min, c.max)
		if c.result != nil {
			assert.Success(err, t)
		}
		assert.EqualWithin(result, c.result, 1e-15, t)
	}
}
