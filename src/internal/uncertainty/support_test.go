package uncertainty

import (
	"testing"

	"github.com/ready-steady/assert"
	"github.com/ready-steady/linear/decomposition"
	"github.com/ready-steady/linear/matrix"
)

func TestInverse(t *testing.T) {
	m := uint(3)

	A := []float64{
		1.0, 2.0, 3.0,
		2.0, 4.0, 5.0,
		3.0, 5.0, 6.0,
	}
	U := make([]float64, m*m)
	Λ := make([]float64, m)

	err := decomposition.SymmetricEigen(A, U, Λ, m)
	assert.Equal(err, nil, t)

	err = matrix.Invert(A, m)
	assert.Equal(err, nil, t)

	I, err := invert(U, Λ, m)
	assert.Equal(err, nil, t)
	assert.Close(A, I, 1e-14, t)
}
