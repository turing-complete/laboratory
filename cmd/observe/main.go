package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"runtime"
	"strconv"

	"github.com/ready-steady/sequence"

	"../internal"
)

var (
	outputFile  = flag.String("o", "", "an output file (required)")
	sampleSeed  = flag.String("s", "", "a seed for generating samples")
	sampleCount = flag.String("n", "", "the number of samples")
)

type Config *internal.AssessmentConfig

func main() {
	internal.Run(command)
}

func command(globalConfig *internal.Config) error {
	globalConfig.Probability.VarThreshold = math.Inf(1)

	config := &globalConfig.Assessment
	if len(*sampleSeed) > 0 {
		if number, err := strconv.ParseInt(*sampleSeed, 0, 64); err != nil {
			return err
		} else {
			config.Seed = number
		}
	}
	if len(*sampleCount) > 0 {
		if number, err := strconv.ParseUint(*sampleCount, 0, 64); err != nil {
			return err
		} else {
			config.Samples = uint(number)
		}
	}

	output, err := internal.Create(*outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	problem, err := internal.NewProblem(globalConfig)
	if err != nil {
		return err
	}

	target, err := internal.NewTarget(problem)
	if err != nil {
		return err
	}

	points, err := generate(config, target)
	if err != nil {
		return err
	}

	ni, no := target.Dimensions()
	ns := uint(len(points)) / ni

	if globalConfig.Verbose {
		fmt.Printf("Evaluating the original model at %d points...\n", ns)
	}

	values := internal.Invoke(target, points, uint(runtime.GOMAXPROCS(0)))

	if globalConfig.Verbose {
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

func generate(config Config, target internal.Target) ([]float64, error) {
	if config.Samples == 0 {
		return nil, errors.New("the number of samples should be positive")
	}

	ni, _ := target.Dimensions()
	sequence := sequence.NewSobol(ni, internal.NewSeed(config.Seed))

	return sequence.Next(config.Samples), nil
}
