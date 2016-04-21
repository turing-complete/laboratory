package uncertainty

import (
	"testing"

	"github.com/ready-steady/assert"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
)

func TestNewEpistemic001(t *testing.T) {
	const (
		nt = 10
		σ  = 0.2
	)

	config, _ := config.New("fixtures/001_010.json")
	system, _ := system.New(&config.System)
	reference := system.ReferenceTime()
	uncertainty, _ := newEpistemic(reference, &config.Uncertainty.Time)

	point := make([]float64, nt)
	value := make([]float64, nt)
	for i := 0; i < nt; i++ {
		α := float64(i) / (nt - 1)
		point[i] = α
		value[i] = (1.0 - σ + 2.0*σ*α) * reference[i]
	}

	assert.EqualWithin(uncertainty.Inverse(point), value, 1e-15, t)
}
