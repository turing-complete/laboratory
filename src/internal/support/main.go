package support

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/ready-steady/linear/matrix"
	"github.com/ready-steady/sequence"
)

var (
	nInfinity = math.Inf(-1)
	pInfinity = math.Inf(1)
)

func Combine(A, x, y []float64, m, n uint) {
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

func Generate(ni, ns uint, seed int64) []float64 {
	return sequence.NewSobol(ni, NewSeed(seed)).Next(ns)
}

func NewSeed(seed int64) int64 {
	if seed < 0 {
		seed = time.Now().Unix()
	}
	return seed
}

func ParseNaturalIndex(line string, min, max uint) ([]uint, error) {
	realIndex, err := ParseRealIndex(line, float64(min), float64(max))
	if err != nil {
		return nil, err
	}

	index := make([]uint, len(realIndex))
	for i := range index {
		index[i] = uint(realIndex[i] + 0.5)
	}

	return index, nil
}

var (
	emptyPattern = regexp.MustCompile(`^(^\[\s*]$)?$`)

	arrayPattern = regexp.MustCompile(`^\[([^:]*)]$`)
	commaPattern = regexp.MustCompile(`\s*,\s*`)

	rangePattern = regexp.MustCompile(`^\[(.*)\]$`)
	colonPattern = regexp.MustCompile(`\s*:\s*`)
)

func ParseRealIndex(line string, min, max float64) ([]float64, error) {
	const (
		ε = 1e-8
	)

	index := make([]float64, 0)

	line = strings.Trim(line, " \t")
	if emptyPattern.MatchString(line) {
		start, step, end := min, 1.0, max
		for start < end || math.Abs(start-end) < ε {
			index = append(index, start)
			start += step
		}
	} else if match := arrayPattern.FindStringSubmatch(line); match != nil {
		for _, chunk := range commaPattern.Split(match[1], -1) {
			number, err := strconv.ParseFloat(chunk, 64)
			if err != nil {
				return nil, err
			}
			index = append(index, number)
		}
	} else if match := rangePattern.FindStringSubmatch(line); match != nil {
		var err error
		var start, step, end float64

		chunks := colonPattern.Split(match[1], -1)

		switch len(chunks) {
		case 2:
			start, err = strconv.ParseFloat(chunks[0], 64)
			if err != nil {
				return nil, err
			}
			step = 1
			end, err = strconv.ParseFloat(chunks[1], 64)
			if err != nil {
				return nil, err
			}
		case 3:
			start, err = strconv.ParseFloat(chunks[0], 64)
			if err != nil {
				return nil, err
			}
			step, err = strconv.ParseFloat(chunks[1], 64)
			if err != nil {
				return nil, err
			}
			end, err = strconv.ParseFloat(chunks[2], 64)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New(fmt.Sprintf("cannot parse the index “%s”", line))
		}

		for start < end || math.Abs(start-end) < ε {
			index = append(index, start)
			start += step
		}
	} else {
		return nil, errors.New(fmt.Sprintf("cannot parse the index “%s”", line))
	}

	for i := range index {
		if math.Abs(index[i]-min) < ε {
			index[i] = min
		}
		if math.Abs(index[i]-max) < ε {
			index[i] = max
		}
		if index[i] < min || index[i] > max {
			return nil, errors.New(fmt.Sprintf("the index “%s” is out of range", line))
		}
	}

	return index, nil
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
