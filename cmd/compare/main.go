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

	lc := config.Assessment.Slices
	if lc == 0 {
		lc = 1
	}
	sc := config.Assessment.Samples
	oc := surrogate.Outputs

	if uint32(len(observed)) != lc*sc*oc || uint32(len(predicted)) != lc*sc*oc {
		return errors.New("an invalid dimensionality of the data")
	}

	εμ, εv, εp := 0.0, 0.0, 0.0

	cut := func(data []float64, i, k uint32) []float64 {
		piece := make([]float64, sc)
		for j := uint32(0); j < sc; j++ {
			piece[j] = data[i*sc*oc+j*oc+k]
		}
		return piece
	}

	// Find the maximal errors across all slices and outputs.
	for i := uint32(0); i < lc; i++ {
		for k := uint32(0); k < oc; k++ {
			observed := cut(observed, i, k)
			predicted := cut(predicted, i, k)

			μ1 := statistics.Mean(observed)
			μ2 := statistics.Mean(predicted)
			if ε := math.Abs((μ1 - μ2) / μ1); ε > εμ {
				εμ = ε
			}

			v1 := statistics.Variance(observed)
			v2 := statistics.Variance(predicted)
			if ε := math.Abs((v1 - v2) / v1); ε > εv {
				εv = ε
			}

			if _, _, ε := test.KolmogorovSmirnov(observed, predicted, 0); ε > εp {
				εp = ε
			}
		}
	}

	fmt.Printf("Inputs: %2d, outputs: %4d, level: %2d, nodes: %7d, μ %.2e, v %.2e, p %.2e\n",
		surrogate.Inputs, surrogate.Outputs, surrogate.Level, surrogate.Nodes, εμ, εv, εp)

	return nil
}
