package main

import (
	"flag"
	"math"

	"github.com/ready-steady/linear/metric"
	"github.com/ready-steady/sort"
	"github.com/ready-steady/statistics/distribution"

	"../internal"
)

const (
	momentCount = 2
	metricCount = 3
)

var (
	referenceFile = flag.String("reference", "", "an output file of `observe` (required)")
	observeFile   = flag.String("observe", "", "an output file of `observe` (required)")
	predictFile   = flag.String("predict", "", "an output file of `predict` (required)")
	outputFile    = flag.String("o", "", "an output file (required)")
)

func main() {
	internal.Run(command)
}

func command(_ *internal.Config) error {
	reference, err := internal.Open(*referenceFile)
	if err != nil {
		return err
	}
	defer reference.Close()

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

	rvalues := []float64{}
	if err := reference.Get("values", &rvalues); err != nil {
		return err
	}

	ovalues := []float64{}
	if err := observe.Get("values", &ovalues); err != nil {
		return err
	}

	pvalues := []float64{}
	if err := predict.Get("values", &pvalues); err != nil {
		return err
	}

	pmoments := []float64{}
	if err := predict.Get("moments", &pmoments); err != nil {
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
		r := slice(rvalues, no, i*momentCount, 1)
		o := cumulate(slice(ovalues, no, i*momentCount, 1), steps)
		p := divide(slice(pvalues, no, i*momentCount, 1), ns)
		m := divide(slice(pmoments, no, i*momentCount, momentCount), ns)

		for j := uint(0); j < ns; j++ {
			εo = append(εo, assess(r, nil, o[j], nil)...)
			εp = append(εp, assess(r, nil, p[j], m[j])...)
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

func assess(data1, moments1, data2, moments2 []float64) []float64 {
	μ1, v1 := computeExpVar(data1, moments1)
	μ2, v2 := computeExpVar(data2, moments2)

	result := make([]float64, metricCount)
	result[0] = math.Abs((μ1 - μ2) / μ1)
	result[1] = math.Abs((v1 - v2) / v1)
	result[2] = computeDistance(data1, data2)

	return result
}

func computeExpVar(data, moments []float64) (float64, float64) {
	var μ float64
	if len(moments) > 0 {
		μ = moments[0]
	} else {
		μ = distribution.Mean(data)
	}

	var v float64
	if len(moments) > 1 {
		v = moments[1] - μ*μ
		if v < 0 {
			v = distribution.Variance(data)
		}
	} else {
		v = distribution.Variance(data)
	}

	return μ, v
}

func computeDistance(data1, data2 []float64) float64 {
	edges := detect(data1, data2)

	cdf1 := distribution.CDF(data1, edges)
	cdf2 := distribution.CDF(data2, edges)

	return metric.Uniform(cdf1, cdf2)
}

func detect(data1, data2 []float64) []float64 {
	n1, n2 := len(data1), len(data2)

	edges := make([]float64, 1+n1+n2+1)
	edges[0] = math.Inf(-1)
	copy(edges[1:], data1)
	copy(edges[1+n1:], data2)
	edges[1+n1+n2] = -edges[0]

	return edges[:sort.Unique(edges)]
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
