package uncertainty

import (
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
	uncertainty, _ := NewAleatory(system, &config.Uncertainty)
	base := uncertainty.(*base)

	for i := 0; i < nt; i++ {
		min, max := (1.0-σ)*reference[i], (1.0+σ)*reference[i]
		assert.EqualWithin(base.marginals[i].Invert(0.0), min, 1e-15, t)
		assert.EqualWithin(base.marginals[i].Invert(1.0), max, 1e-15, t)
	}
}

func TestNewEpistemic(t *testing.T) {
	const (
		nt = 10
		σ  = 0.2
	)

	config, _ := config.New("fixtures/001_010.json")
	system, _ := system.New(&config.System)
	reference := system.ReferenceTime()
	uncertainty, _ := NewEpistemic(system, &config.Uncertainty)

	point := make([]float64, nt)
	value := make([]float64, nt)
	for i := 0; i < nt; i++ {
		α := float64(i) / (nt - 1)
		point[i] = α
		value[i] = (1.0 - σ + 2.0*σ*α) * reference[i]
	}

	assert.EqualWithin(uncertainty.Inverse(point), value, 1e-15, t)
}
