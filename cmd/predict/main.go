package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"

	"github.com/ready-steady/sequence"

	"../internal"
)

var (
	approximateFile = flag.String("approximate", "", "an output of `approximate` (required)")
	outputFile      = flag.String("o", "", "an output file (required)")
	sampleSeed      = flag.String("s", "", "a seed for generating samples")
	sampleCount     = flag.String("n", "", "the number of samples")
)

type Config *internal.AssessmentConfig

func main() {
	internal.Run(command)
}

func command(globalConfig *internal.Config) error {
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

	approximate, err := internal.Open(*approximateFile)
	if err != nil {
		return err
	}
	defer approximate.Close()

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

	solver, err := internal.NewSolver(problem, target)
	if err != nil {
		return err
	}

	solution := new(internal.Solution)
	if err = approximate.Get("solution", solution); err != nil {
		return err
	}

	points, err := generate(config, target)
	if err != nil {
		return err
	}

	ni, no := target.Dimensions()
	ns := uint(len(points)) / ni
	np := uint(len(solution.Steps))

	if globalConfig.Verbose {
		fmt.Printf("Evaluating the surrogate model %d times at %d points...\n", np, ns)
		fmt.Printf("%10s %15s\n", "Step", "Nodes")
	}

	values := make([]float64, np*ns*no)
	moments := make([]float64, np*no)

	for i, nn := uint(0), uint(0); i < np; i++ {
		nn += solution.Steps[i]

		if globalConfig.Verbose {
			fmt.Printf("%10d %15d\n", i, nn)
		}

		s := *solution
		s.Nodes = nn
		s.Indices = s.Indices[:nn*ni]
		s.Surpluses = s.Surpluses[:nn*no]

		copy(values[i*ns*no:(i+1)*ns*no], solver.Evaluate(&s, points))
		copy(moments[i*no:(i+1)*no], solver.Integrate(&s))
	}

	if globalConfig.Verbose {
		fmt.Println("Done.")
	}

	if err := output.Put("solution", *solution); err != nil {
		return err
	}
	if err := output.Put("points", points, ni, ns); err != nil {
		return err
	}
	if err := output.Put("values", values, no, ns, np); err != nil {
		return err
	}
	if err := output.Put("moments", moments, no, np); err != nil {
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
