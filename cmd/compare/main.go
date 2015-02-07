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

	observed := []float64{}
	if err := input.Get("observed", &observed); err != nil {
		return err
	}

	predicted := []float64{}
	if err := input.Get("predicted", &predicted); err != nil {
		return err
	}

	μ1 := statistics.Mean(observed)
	μ2 := statistics.Mean(predicted)

	σ1 := math.Sqrt(statistics.Variance(observed))
	σ2 := math.Sqrt(statistics.Variance(predicted))

	_, _, Δ := test.KolmogorovSmirnov(observed, predicted, 0)

	fmt.Printf("Inputs: %2d, outputs: %4d, level: %2d, nodes: %7d, μ %.2e, σ %.2e, Δ %.2e\n",
		surrogate.Inputs, surrogate.Outputs, surrogate.Level, surrogate.Nodes,
		math.Abs((μ1 - μ2) / μ1), math.Abs((σ1 - σ2) / σ1), Δ)

	return nil
}
