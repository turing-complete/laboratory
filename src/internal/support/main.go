package support

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ready-steady/sequence"
)

var (
	emptyPattern = regexp.MustCompile(`^\[\s*]$`)

	arrayPattern = regexp.MustCompile(`^\[([^:]*)]$`)
	commaPattern = regexp.MustCompile(`\s*,\s*`)

	rangePattern = regexp.MustCompile(`^\[(.*)\]$`)
	colonPattern = regexp.MustCompile(`\s*:\s*`)
)

func Average(data []float64) float64 {
	return Sum(data) / float64(len(data))
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

func ParseRealIndex(line string, min, max float64) ([]float64, error) {
	const (
		ε = 1e-8
	)

	index := make([]float64, 0)

	parse := func(chunk string) (float64, error) {
		if chunk == "end" {
			return max, nil
		} else {
			number, err := strconv.ParseFloat(chunk, 64)
			if err != nil {
				return 0.0, err
			}
			return number, nil
		}
	}

	line = strings.ToLower(strings.Trim(line, " \t"))
	if emptyPattern.MatchString(line) {
	} else if match := arrayPattern.FindStringSubmatch(line); match != nil {
		for _, chunk := range commaPattern.Split(match[1], -1) {
			number, err := parse(chunk)
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
			start, err = parse(chunks[0])
			if err != nil {
				return nil, err
			}
			step = 1
			end, err = parse(chunks[1])
			if err != nil {
				return nil, err
			}
		case 3:
			start, err = parse(chunks[0])
			if err != nil {
				return nil, err
			}
			step, err = parse(chunks[1])
			if err != nil {
				return nil, err
			}
			end, err = parse(chunks[2])
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

func Sum(data []float64) (Σ float64) {
	for _, x := range data {
		Σ += x
	}
	return
}
