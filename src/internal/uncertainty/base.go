package uncertainty

import (
	"errors"
	"fmt"
	"math"

	"github.com/ready-steady/linear/matrix"
	"github.com/ready-steady/probability"
	"github.com/ready-steady/statistics/correlation"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/support"
	"github.com/turing-complete/laboratory/src/internal/system"

	icorrelation "github.com/turing-complete/laboratory/src/internal/correlation"
)

var (
	standardGaussian = probability.NewGaussian(0, 1)
	nInfinity        = math.Inf(-1)
	pInfinity        = math.Inf(1)
)

type base struct {
	taskIndex  []uint
	correlator []float64

	nt uint
	nu uint
	nz uint
}

func newBase(c *config.Uncertainty, s *system.System) (*base, error) {
	nt := uint(s.Application.Len())

	taskIndex, err := support.ParseNaturalIndex(c.TaskIndex, 0, nt-1)
	if err != nil {
		return nil, err
	}

	correlator, err := correlate(c, s, taskIndex)
	if err != nil {
		return nil, err
	}

	nu := uint(len(taskIndex))
	nz := uint(len(correlator)) / nu

	base := &base{
		taskIndex:  taskIndex,
		correlator: correlator,

		nt: nt,
		nu: nu,
		nz: nz,
	}

	return base, nil
}

func (b *base) Len() int {
	return int(b.nz)
}

func (b *base) String() string {
	return fmt.Sprintf(`{"parameters": %d, "variables": %d}`, b.nu, b.nz)
}

func (b *base) Transform(z []float64) []float64 {
	nu, nz := b.nu, b.nz

	n := make([]float64, nz)
	u := make([]float64, nu)

	// Independent uniform to independent Gaussian
	for i := range n {
		n[i] = standardGaussian.InvCDF(z[i])
	}

	// Independent Gaussian to dependent Gaussian
	multiply(b.correlator, n, u, nu, nz)

	// Dependent Gaussian to dependent uniform
	for i := range u {
		u[i] = standardGaussian.CDF(u[i])
	}

	return u
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
