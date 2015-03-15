package internal

import (
	"errors"
	"fmt"
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
