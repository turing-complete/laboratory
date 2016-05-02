package uncertainty

import (
	"errors"
	"math"

	"github.com/ready-steady/linear/matrix"
)

var (
	infinity = math.Inf(1.0)
)

func inverse(U, Λ []float64, m uint) ([]float64, error) {
	T := make([]float64, m*m)
	for i := uint(0); i < m; i++ {
		if Λ[i] == 0.0 {
			return nil, errors.New("the matrix is not invertible")
		}
		λ := 1.0 / Λ[i]
		for j := uint(0); j < m; j++ {
			T[j*m+i] = λ * U[i*m+j]
		}
	}

	I := make([]float64, m*m)
	matrix.Multiply(U, T, I, m, m, m)

	return I, nil
}

func multiply(A, x, y []float64, m, n uint) {
	infinite, z := false, make([]float64, n)

	for i := range x {
		switch x[i] {
		case -infinity:
			infinite, z[i] = true, -1.0
		case infinity:
			infinite, z[i] = true, 1.0
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
			if a == 0.0 {
				continue
			}
			if z[j] == 0.0 {
				Σ1 += a * x[j]
			} else {
				Σ2 += a * z[j]
			}
		}
		if Σ2 < 0.0 {
			y[i] = -infinity
		} else if Σ2 > 0.0 {
			y[i] = infinity
		} else {
			y[i] = Σ1
		}
	}
}
