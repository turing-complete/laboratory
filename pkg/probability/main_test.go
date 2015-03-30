package probability

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestParseInverter(t *testing.T) {
	cases := []struct {
		line    string
		success bool
	}{
		{"Beta(1, 1)", true},
		{"beta(0.5, 1.5)", true},
		{" Beta \t (1, 1)", true},
		{"Gamma(1, 1)", false},
		{"Beta(1, 1, 1)", false},
		{"beta(-1, 1)", false},
		{"beta(0, 1)", false},
		{"beta(1, -1)", false},
		{"beta(1, 0)", false},
		{"beta(1, 0)", false},
	}

	for _, c := range cases {
		if _, err := ParseInverter(c.line); c.success {
			assert.Success(err, t)
		} else {
			assert.Failure(err, t)
		}
	}
}
