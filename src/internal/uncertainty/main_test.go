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
	uncertainty, _ := NewAleatory(system, &config.Uncertainty)

	point := make([]float64, nt)
	for i := 0; i < nt; i++ {
		point[i] = 0.5
	}

	assert.EqualWithin(uncertainty.Backward(point), []float64{
		3.1402438661763954e-02, 1.7325483399593899e-02,
		2.7071067811865485e-02, 3.1402438661763954e-02,
		4.0065180361560912e-02, 3.2485281374238568e-02,
		1.7325483399593888e-02, 2.5988225099390850e-02,
		1.6242640687119302e-02, 3.2485281374238568e-02,
	}, 1e-15, t)
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

	assert.EqualWithin(uncertainty.Backward(point), value, 1e-15, t)
}
