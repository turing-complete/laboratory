package main

import (
	"errors"
	"fmt"
	"math"

	"github.com/ready-steady/format/mat"
	"github.com/ready-steady/numeric/interpolation/adhier"
	"github.com/ready-steady/statistics"
	"github.com/ready-steady/statistics/test"

	"../internal"
)

const (
	deltaCensiusKelvin = 273.15
)

func main() {
	internal.Run(command)
}

func command(config internal.Config, input *mat.File, _ *mat.File) error {
	if input == nil {
		return errors.New("an input file is required")
	}

	surrogate := new(adhier.Surrogate)
	if err := input.Get("surrogate", surrogate); err != nil {
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

	nt := config.Assessment.Steps
	if nt == 0 {
		nt = 1
	}
	ns := config.Assessment.Samples
	no := uint(len(observations)) / (nt * ns)

	cut := func(data []float64, i, j uint) []float64 {
		piece := make([]float64, ns)
		for k := uint(0); k < ns; k++ {
			piece[k] = data[i*ns*no+k*no+j]
		}
		return piece
	}

	fmt.Printf("Surrogate: inputs %d, outputs %d, level %d, nodes %d\n",
		surrogate.Inputs, surrogate.Outputs, surrogate.Level, surrogate.Nodes)

	εμ := make([]float64, nt*no)
	εσ := make([]float64, nt*no)
	εp := make([]float64, nt*no)

	// Compute errors across all time moments and outputs.
	for i := uint(0); i < nt; i++ {
		for j := uint(0); j < no; j++ {
			k := i*no + j

			observations := cut(observations, i, j)
			predictions := cut(predictions, i, j)

			μ1 := statistics.Mean(observations)
			μ2 := statistics.Mean(predictions)
			εμ[k] = math.Abs(μ1 - μ2)

			σ1 := math.Sqrt(statistics.Variance(observations))
			σ2 := math.Sqrt(statistics.Variance(predictions))
			εσ[k] = math.Abs(σ1 - σ2)

			_, _, εp[k] = test.KolmogorovSmirnov(observations, predictions, 0)

			if config.Verbose {
				fmt.Printf("%9d: μ %10.4f ±%10.4f, σ %10.4f ±%10.4f, p %.2e\n",
					k, μ1-deltaCensiusKelvin, εμ[k], σ1, εσ[k], εp[k])
			}
		}
	}

	fmt.Printf("Average error: μ ±%10.4f, σ ±%10.4f, p %.2e\n",
		statistics.Mean(εμ), statistics.Mean(εσ), statistics.Mean(εp))

	fmt.Printf("Maximal error: μ ±%10.4f, σ ±%10.4f, p %.2e\n",
		max(εμ), max(εσ), max(εp))

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
