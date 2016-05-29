package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/ready-steady/statistics/metric"
	"github.com/turing-complete/laboratory/src/internal/command"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/database"
	"github.com/turing-complete/laboratory/src/internal/solution"
)

const (
	momentCount = 1
	metricCount = 1
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

	pvalues := []float64{}
	if err := predict.Get("values", &pvalues); err != nil {
		return err
	}

	active := []uint{}
	if err := predict.Get("active", &active); err != nil {
		return err
	}

	surrogate := new(solution.Surrogate)
	if err := predict.Get("surrogate", surrogate); err != nil {
		return err
	}

	no := surrogate.Outputs
	nq := no / momentCount
	nk := uint(len(active))

	if ne := active[nk-1]; uint(len(ovalues))/no < ne {
		return errors.New(fmt.Sprintf("the number of observations should be at least %d", ne))
	}

	εo := make([]float64, 0, nq*nk*metricCount)
	εp := make([]float64, 0, nq*nk*metricCount)

	for i := uint(0); i < nq; i++ {
		r := slice(rvalues, no, i*momentCount, 1)

		o := cumulate(slice(ovalues, no, i*momentCount, 1), active)
		for j := uint(0); j < nk; j++ {
			εo = append(εo, assess(r, o[j])...)
		}

		p := divide(slice(pvalues, no, i*momentCount, 1), nk)
		for j := uint(0); j < nk; j++ {
			εp = append(εp, assess(r, p[j])...)
		}
	}

	if err := output.Put("active", active); err != nil {
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
	return []float64{metric.KolmogorovSmirnov(data1, data2)}
}

func cumulate(data []float64, cumsum []uint) [][]float64 {
	n := uint(len(cumsum))
	sets := make([][]float64, n)
	for i := uint(0); i < n; i++ {
		sets[i] = data[:cumsum[i]]
	}
	return sets
}

func divide(data []float64, n uint) [][]float64 {
	step := uint(len(data)) / n
	sets := make([][]float64, n)
	for i := uint(0); i < n; i++ {
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
