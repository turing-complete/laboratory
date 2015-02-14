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

	tc := config.Assessment.Times
	if tc == 0 {
		tc = 1
	}
	sc := config.Assessment.Samples
	oc := uint(len(observations)) / (tc * sc)

	εμ, εv, εp := 0.0, 0.0, 0.0

	cut := func(data []float64, i, k uint) []float64 {
		piece := make([]float64, sc)
		for j := uint(0); j < sc; j++ {
			piece[j] = data[i*sc*oc+j*oc+k]
		}
		return piece
	}

	// Find the maximal errors across all time moments and outputs.
	for i := uint(0); i < tc; i++ {
		for k := uint(0); k < oc; k++ {
			observations := cut(observations, i, k)
			predictions := cut(predictions, i, k)

			μ1 := statistics.Mean(observations)
			μ2 := statistics.Mean(predictions)
			if ε := math.Abs((μ1 - μ2) / μ1); ε > εμ {
				εμ = ε
			}

			v1 := statistics.Variance(observations)
			v2 := statistics.Variance(predictions)
			if ε := math.Abs((v1 - v2) / v1); ε > εv {
				εv = ε
			}

			if _, _, ε := test.KolmogorovSmirnov(observations, predictions, 0); ε > εp {
				εp = ε
			}
		}
	}

	fmt.Printf("Inputs: %2d, outputs: %4d, level: %2d, nodes: %7d, μ %.2e, v %.2e, p %.2e\n",
		surrogate.Inputs, surrogate.Outputs, surrogate.Level, surrogate.Nodes, εμ, εv, εp)

	return nil
}
