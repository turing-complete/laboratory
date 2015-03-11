package main

import (
	"errors"
	"fmt"
	"math"

	"github.com/ready-steady/format/mat"
	"github.com/ready-steady/statistics"
	"github.com/ready-steady/statistics/test"

	"../internal"
)

func main() {
	internal.Run(command)
}

func command(config internal.Config, input *mat.File, _ *mat.File) error {
	if input == nil {
		return errors.New("an input file is required")
	}

	solution := new(internal.Solution)
	if err := input.Get("solution", solution); err != nil {
		return err
	}

	observations := []float64{}
	if err := input.Get("observations", &observations); err != nil {
		return err
	}

	predictions := []float64{}
	if err := input.Get("predictions", &predictions); err != nil {
		return err
	}

	ns := int(config.Assessment.Samples)
	no := len(observations) / ns
	nm := no / 2

	cut := func(data []float64, i int) []float64 {
		piece := make([]float64, ns)
		for j := 0; j < ns; j++ {
			piece[j] = data[j*no+i]
		}
		return piece
	}

	fmt.Println(solution)

	εμ := make([]float64, nm)
	εv := make([]float64, nm)
	εp := make([]float64, nm)

	// Compute errors across all outputs.
	for i := 0; i < nm; i++ {
		j := i * 2

		observations := cut(observations, j)
		predictions := cut(predictions, j)

		μ1 := statistics.Mean(observations)
		μ2 := solution.Expectation[j]
		εμ[i] = math.Abs(μ1 - μ2)

		v1 := statistics.Variance(observations)
		v2 := solution.Expectation[j+1] - μ2*μ2
		εv[i] = math.Abs(v1 - v2)

		_, _, εp[i] = test.KolmogorovSmirnov(observations, predictions, 0)

		if nm == 1 {
			fmt.Printf("Error: μ %.2e ± %.2e (%.2e), v %.2e ± %.2e (%.2e), p %.2e\n",
				μ1, εμ[i], εμ[i]/μ1, v1, εv[i], εv[i]/v1, εp[i])
		} else if config.Verbose {
			fmt.Printf("%9d: μ %.2e ± %.2e (%.2e), v %.2e ± %.2e (%.2e), p %.2e\n",
				i, μ1, εμ[i], εμ[i]/μ1, v1, εv[i], εv[i]/v1, εp[i])
		}
	}

	if nm > 1 {
		fmt.Printf("Average error: μ ± %.2e, v ± %.2e, p %.2e\n",
			statistics.Mean(εμ), statistics.Mean(εv), statistics.Mean(εp))

		fmt.Printf("Maximal error: μ ± %.2e, v ± %.2e, p %.2e\n",
			max(εμ), max(εv), max(εp))
	}

	return nil
}

func max(data []float64) float64 {
	max := math.Inf(-1)

	for _, x := range data {
		if x > max {
			max = x
		}
	}

	return max
}
