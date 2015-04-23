package main

import (
	"flag"
	"math"

	"github.com/ready-steady/statistics"
	"github.com/ready-steady/statistics/test"

	"../internal"
)

const (
	momentCount = 2
	metricCount = 3
)

var (
	observeFile = flag.String("observe", "", "an output file of `observe` (required)")
	predictFile = flag.String("predict", "", "an output file of `predict` (required)")
	outputFile  = flag.String("o", "", "an output file (required)")
)

func main() {
	internal.Run(command)
}

func command(_ *internal.Config) error {
	observe, err := internal.Open(*observeFile)
	if err != nil {
		return err
	}
	defer observe.Close()

	predict, err := internal.Open(*predictFile)
	if err != nil {
		return err
	}
	defer predict.Close()

	output, err := internal.Create(*outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	observations := []float64{}
	if err := observe.Get("values", &observations); err != nil {
		return err
	}

	predictions := []float64{}
	if err := predict.Get("values", &predictions); err != nil {
		return err
	}

	moments := []float64{}
	if err := predict.Get("moments", &moments); err != nil {
		return err
	}

	solution := new(internal.Solution)
	if err := predict.Get("solution", solution); err != nil {
		return err
	}

	steps := solution.Steps

	no := solution.Outputs
	nq := no / momentCount
	ns := uint(len(steps))

	εo := make([]float64, 0, nq*ns*metricCount)
	εp := make([]float64, 0, nq*ns*metricCount)

	for i := uint(0); i < nq; i++ {
		observations := slice(observations, no, i*momentCount, 1)
		predictions := slice(predictions, no, i*momentCount, 1)
		moments := slice(moments, no, i*momentCount, momentCount)

		data1 := cumulate(observations, steps)
		data2 := divide(predictions, ns)
		mean2 := divide(moments, ns)

		for j := uint(0); j < ns; j++ {
			εo = append(εo, compare(observations, data1[j], nil)...)
			εp = append(εp, compare(observations, data2[j], mean2[j])...)
		}
	}

	if err := output.Put("steps", steps); err != nil {
		return err
	}

	if err := output.Put("observe", εo, metricCount, ns, nq); err != nil {
		return err
	}

	if err := output.Put("predict", εp, metricCount, ns, nq); err != nil {
		return err
	}

	return nil
}

func compare(data1, data2, mean2 []float64) []float64 {
	μ1 := statistics.Mean(data1)
	v1 := statistics.Variance(data1)

	var μ2 float64
	if len(mean2) > 0 {
		μ2 = mean2[0]
	} else {
		μ2 = statistics.Mean(data2)
	}

	var v2 float64
	if len(mean2) > 1 {
		v2 = mean2[1] - μ2*μ2
		if v2 < 0 {
			v2 = statistics.Variance(data2)
		}
	} else {
		v2 = statistics.Variance(data2)
	}

	εμ := math.Abs((μ1 - μ2) / μ1)
	εv := math.Abs((v1 - v2) / v1)
	_, _, εp := test.KolmogorovSmirnov(data1, data2, 0)

	result := make([]float64, metricCount)
	result[0] = εμ
	result[1] = εv
	result[2] = εp

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
