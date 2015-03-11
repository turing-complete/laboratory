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

	cut := func(data []float64, i int) []float64 {
		piece := make([]float64, ns)
		for j := 0; j < ns; j++ {
			piece[j] = data[j*no+i]
		}
		return piece
	}

	fmt.Println(solution)

	εμ := make([]float64, no)
	εσ := make([]float64, no)
	εp := make([]float64, no)

	// Compute errors across all outputs.
	for i := 0; i < no; i++ {
		observations := cut(observations, i)
		predictions := cut(predictions, i)

		μ1 := statistics.Mean(observations)
		μ2 := statistics.Mean(predictions)
		εμ[i] = math.Abs(μ1 - μ2)

		σ1 := math.Sqrt(statistics.Variance(observations))
		σ2 := math.Sqrt(statistics.Variance(predictions))
		εσ[i] = math.Abs(σ1 - σ2)

		_, _, εp[i] = test.KolmogorovSmirnov(observations, predictions, 0)

		if no == 1 {
			fmt.Printf("Error: μ %10.2e ±%10.2e, σ %10.2e ±%10.2e, p %.2e\n",
				μ1, εμ[i], σ1, εσ[i], εp[i])
		} else if config.Verbose {
			fmt.Printf("%9d: μ %10.2e ±%10.2e, σ %10.2e ±%10.2e, p %.2e\n",
				i, μ1, εμ[i], σ1, εσ[i], εp[i])
		}
	}

	if no > 1 {
		fmt.Printf("Average error: μ ±%10.2e, σ ±%10.2e, p %.2e\n",
			statistics.Mean(εμ), statistics.Mean(εσ), statistics.Mean(εp))

		fmt.Printf("Maximal error: μ ±%10.2e, σ ±%10.2e, p %.2e\n",
			max(εμ), max(εσ), max(εp))
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
