package uncertainty

import (
	"errors"

	"github.com/ready-steady/linear/matrix"
)

func invert(U, Λ []float64, m uint) ([]float64, error) {
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
