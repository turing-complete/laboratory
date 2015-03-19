package internal

import (
	"errors"
	"math"
	"reflect"
	"unsafe"

	"github.com/ready-steady/linear/matrix"
)

var (
	nInfinity = math.Inf(-1)
	pInfinity = math.Inf(1)
)

func combine(A, x, y []float64, m, n uint) {
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

func enumerate(count uint, line []uint) ([]uint, error) {
	if len(line) == 0 {
		line = make([]uint, count)
		for i := uint(0); i < count; i++ {
			line[i] = i
		}
	}

	for _, i := range line {
		if i >= count {
			return nil, errors.New("the index is out of range")
		}
	}

	return line, nil
}

func locate(l, r float64, line []float64) (uint, uint) {
	n := len(line)

	i, j := 0, n-1

	for i < j-1 {
		if l < line[i+1] {
			break
		}
		i++
	}
	for j > i+1 {
		if r > line[j-1] {
			break
		}
		j--
	}

	return uint(i), uint(j + 1)
}

func slice(data []float64, index []uint, m uint) []float64 {
	n := uint(len(data)) / m
	p := uint(len(index))

	chunk := make([]float64, p*n)

	for i := uint(0); i < p; i++ {
		for j := uint(0); j < n; j++ {
			chunk[j*p+i] = data[j*m+index[i]]
		}
	}

	return chunk
}

func stringify(node []float64) string {
	const (
		sizeOfFloat64 = 8
	)

	var bytes []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	header.Data = ((*reflect.SliceHeader)(unsafe.Pointer(&node))).Data
	header.Cap = sizeOfFloat64 * len(node)
	header.Len = header.Cap

	return string(bytes)
}

func subdivide(span, Δx float64, fraction []float64) ([]float64, error) {
	if Δx <= 0 {
		return nil, errors.New("the step should be positive")
	}

	var left, right float64

	switch len(fraction) {
	case 0:
		left, right = 0, span
	case 1:
		left, right = fraction[0]*span, fraction[0]*span
	default:
		left, right = fraction[0]*span, fraction[1]*span
	}
	if left < 0 || left > right || right > span {
		return nil, errors.New("the fraction should be between 0 and 1")
	}

	line := make([]float64, 0, uint((right-left)/Δx)+1)
	for t := left; t <= right; t += Δx {
		line = append(line, t)
	}

	return line, nil
}
