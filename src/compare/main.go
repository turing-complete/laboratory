package main

import (
	"errors"
	"flag"
	"fmt"
	"math"

	"github.com/ready-steady/statistics/distribution"
	"github.com/ready-steady/statistics/metric"
	"github.com/turing-complete/laboratory/src/internal/command"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/database"
	"github.com/turing-complete/laboratory/src/internal/solver"
)

const (
	momentCount = 1
	metricCount = 3
)

var (
	referenceFile = flag.String("reference", "", "an output file of `observe` (required)")
	observeFile   = flag.String("observe", "", "an output file of `observe` (required)")
	predictFile   = flag.String("predict", "", "an output file of `predict` (required)")
	outputFile    = flag.String("o", "", "an output file (required)")
)

func main() {
	command.Run(function)
}

func function(_ *config.Config) error {
	reference, err := database.Open(*referenceFile)
	if err != nil {
		return err
	}
	defer reference.Close()

	observe, err := database.Open(*observeFile)
	if err != nil {
		return err
	}
	defer observe.Close()

	predict, err := database.Open(*predictFile)
	if err != nil {
		return err
	}
	defer predict.Close()

	output, err := database.Create(*outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	rvalues := []float64{}
	if err := reference.Get("values", &rvalues); err != nil {
		return err
	}

	ovalues := []float64{}
	if err := observe.Get("values", &ovalues); err != nil {
		return err
	}

	psteps := []uint{}
	if err := predict.Get("steps", &psteps); err != nil {
		return err
	}

	pvalues := []float64{}
	if err := predict.Get("values", &pvalues); err != nil {
		return err
	}

	solution := new(solver.Solution)
	if err := predict.Get("solution", solution); err != nil {
		return err
	}

	no := solution.Outputs
	nq := no / momentCount
	nk := uint(len(psteps))

	ne := uint(0)
	for _, step := range psteps {
		ne += step
	}
	if uint(len(ovalues))/no < ne {
		return errors.New(fmt.Sprintf("the number of observations should be at least %d", ne))
	}

	εo := make([]float64, 0, nq*nk*metricCount)
	εp := make([]float64, 0, nq*nk*metricCount)

	for i := uint(0); i < nq; i++ {
		r := slice(rvalues, no, i*momentCount, 1)

		o := cumulate(slice(ovalues, no, i*momentCount, 1), psteps)
		for j := uint(0); j < nk; j++ {
			εo = append(εo, assess(r, o[j])...)
		}

		p := divide(slice(pvalues, no, i*momentCount, 1), nk)
		for j := uint(0); j < nk; j++ {
			εp = append(εp, assess(r, p[j])...)
		}
	}

	if err := output.Put("steps", psteps); err != nil {
		return err
	}
	if err := output.Put("observe", εo, metricCount, nk, nq); err != nil {
		return err
	}
	if err := output.Put("predict", εp, metricCount, nk, nq); err != nil {
		return err
	}

	return nil
}

func assess(data1, data2 []float64) []float64 {
	μ1, v1 := distribution.Expectation(data1), distribution.Variance(data1)
	μ2, v2 := distribution.Expectation(data2), distribution.Variance(data2)

	result := make([]float64, metricCount)
	result[0] = math.Abs((μ1 - μ2) / μ1)
	result[1] = math.Abs((v1 - v2) / v1)
	result[2] = metric.KolmogorovSmirnov(data1, data2)

	return result
}

func cumulate(data []float64, steps []uint) [][]float64 {
	count := uint(len(steps))

	sets := make([][]float64, count)
	for i, sum := uint(0), uint(0); i < count; i++ {
		sum += steps[i]
		sets[i] = data[:sum]
	}

	return sets
}

func divide(data []float64, count uint) [][]float64 {
	step := uint(len(data)) / count

	sets := make([][]float64, count)
	for i := uint(0); i < count; i++ {
		sets[i] = data[i*step : (i+1)*step]
	}

	return sets
}

func slice(data []float64, height, offset, thickness uint) []float64 {
	width := uint(len(data)) / height
	piece := make([]float64, thickness*width)

	for i := uint(0); i < thickness; i++ {
		for j := uint(0); j < width; j++ {
			piece[j*thickness+i] = data[j*height+offset+i]
		}
	}

	return piece
}
