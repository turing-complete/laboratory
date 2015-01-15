package probability

import (
	"testing"
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
		result := ParseInverter(c.line)
		if c.success && result == nil {
			t.Errorf("expected “%v” to succeed", c.line)
		} else if !c.success && result != nil {
			t.Errorf("expected “%v” to fail", c.line)
		}
	}
}
