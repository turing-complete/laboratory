package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"runtime"

	"github.com/ready-steady/sequence"

	"../internal"
)

var (
	outputFile  = flag.String("o", "", "an output file (required)")
	sampleSeed  = flag.Float64("s", math.NaN(), "a seed for generating samples")
	sampleCount = flag.Float64("n", math.NaN(), "the number of samples")
)

func main() {
	internal.Run(command)
}

func command(config *internal.Config) error {
	output, err := internal.Create(*outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	config.Probability.VarThreshold = math.Inf(1)
	if !math.IsNaN(*sampleSeed) {
		config.Assessment.Seed = int64(*sampleSeed)
	}
	if !math.IsNaN(*sampleCount) {
		config.Assessment.Samples = uint(*sampleCount)
	}

	problem, err := internal.NewProblem(config)
	if err != nil {
		return err
	}

	target, err := internal.NewTarget(problem)
	if err != nil {
		return err
	}

	points, err := generate(&config.Assessment, target)
	if err != nil {
		return err
	}

	ni, no := target.Dimensions()
	ns := uint(len(points)) / ni

	if config.Verbose {
		fmt.Printf("Evaluating the original model at %d points...\n", ns)
	}

	values := internal.Invoke(target, points, uint(runtime.GOMAXPROCS(0)))

	if config.Verbose {
		fmt.Println("Done.")
	}

	if err := output.Put("points", points, ni, ns); err != nil {
		return err
	}
	if err := output.Put("values", values, no, ns); err != nil {
		return err
	}

	return nil
}

func generate(config *internal.AssessmentConfig, target internal.Target) ([]float64, error) {
	if config.Samples == 0 {
		return nil, errors.New("the number of samples should be positive")
	}

	ni, _ := target.Dimensions()
	sequence := sequence.NewSobol(ni, internal.NewSeed(config.Seed))

	return sequence.Next(config.Samples), nil
}
