package internal

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

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

func subdivide(interval []float64, Δx, span float64) ([]float64, error) {
	if Δx <= 0 {
		return nil, errors.New("the step should be positive")
	}

	var left, right float64

	switch len(interval) {
	case 0:
		left, right = 0, span
	case 1:
		left, right = interval[0], interval[0]
	default:
		left, right = interval[0], interval[1]
	}
	if left < 0 || left > right || right > span {
		return nil, errors.New(fmt.Sprintf("the interval should be between 0 and %g", span))
	}

	line := make([]float64, 0, uint((right-left)/Δx)+1)
	for t := left; t <= right; t += Δx {
		line = append(line, t)
	}

	return line, nil
}
