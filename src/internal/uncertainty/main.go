package uncertainty

import (
	"errors"
	"math"

	"github.com/ready-steady/linear/matrix"
	"github.com/ready-steady/statistics/correlation"
	"github.com/simulated-reality/laboratory/src/internal/config"
	"github.com/simulated-reality/laboratory/src/internal/system"

	icorrelation "github.com/simulated-reality/laboratory/src/internal/correlation"
)

var (
	nInfinity = math.Inf(-1)
	pInfinity = math.Inf(1)
)

type Uncertainty interface {
	Len() int
	Transform([]float64) []float64
}

func correlate(c *config.Uncertainty, s *system.System, taskIndex []uint) ([]float64, error) {
	if c.CorrLength < 0 {
		return nil, errors.New("the correlation length should be nonnegative")
	}
	if c.VarThreshold <= 0 {
		return nil, errors.New("the variance-reduction threshold should be positive")
	}

	C := icorrelation.Compute(s.Application, taskIndex, c.CorrLength)
	correlator, _, err := correlation.Decompose(C, uint(len(taskIndex)), c.VarThreshold)
	if err != nil {
		return nil, err
	}

	return correlator, nil
}

func multiply(A, x, y []float64, m, n uint) {
	infinite, z := false, make([]float64, n)

	for i := range x {
		switch x[i] {
		case nInfinity:
			infinite, z[i] = true, -1
		case pInfinity:
			infinite, z[i] = true, 1
		}
	}

	if !infinite {
		matrix.Multiply(A, x, y, m, n, 1)
		return
	}

	for i := uint(0); i < m; i++ {
		Σ1, Σ2 := 0.0, 0.0
		for j := uint(0); j < n; j++ {
			a := A[j*m+i]
			if a == 0 {
				continue
			}
			if z[j] == 0 {
				Σ1 += a * x[j]
			} else {
				Σ2 += a * z[j]
			}
		}
		if Σ2 < 0 {
			y[i] = nInfinity
		} else if Σ2 > 0 {
			y[i] = pInfinity
		} else {
			y[i] = Σ1
		}
	}
}
