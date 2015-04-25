package main

import (
	"errors"
	"flag"
	"fmt"
	"runtime"

	"github.com/ready-steady/sequence"

	"../internal"
)

var (
	outputFile = flag.String("o", "", "an output file (required)")
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

	config.Probability.VarThreshold = 42

	problem, err := internal.NewProblem(config)
	if err != nil {
		return err
	}

	target, err := internal.NewTarget(problem)
	if err != nil {
		return err
	}

	points, err := generate(problem, target)
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

func generate(problem *internal.Problem, target internal.Target) ([]float64, error) {
	config := &problem.Config.Assessment
	if config.ReferenceSamples == 0 {
		return nil, errors.New("the number of samples should be positive")
	}

	ni, _ := target.Dimensions()
	sequence := sequence.NewSobol(ni, internal.NewSeed(config.Seed))

	return sequence.Next(config.ReferenceSamples), nil
}
