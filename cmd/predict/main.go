package main

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/ready-steady/sequence"

	"../internal"
)

var (
	constructFile = flag.String("construct", "", "an output of `construct` (required)")
	outputFile    = flag.String("o", "", "an output file (required)")
)

func main() {
	internal.Run(command)
}

func command(config *internal.Config) error {
	construct, err := internal.Open(*constructFile)
	if err != nil {
		return err
	}
	defer construct.Close()

	output, err := internal.Create(*outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	problem, err := internal.NewProblem(config)
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
	if err = construct.Get("solution", solution); err != nil {
		return err
	}

	points, err := generate(problem, target)
	if err != nil {
		return err
	}

	if config.Verbose {
		fmt.Println("Sampling the surrogate model...")
	}

	ni, no := target.Dimensions()
	ns := config.Assessment.Samples
	np := uint(len(solution.Steps))

	values := make([]float64, np*ns*no)
	moments := make([]float64, np*no)

	for i, nn := uint(0), uint(0); i < np; i++ {
		nn += solution.Steps[i]

		if config.Verbose {
			fmt.Printf("%5d: %10d\n", i, nn)
		}

		s := *solution
		s.Nodes = nn
		s.Indices = s.Indices[:nn*ni]
		s.Surpluses = s.Surpluses[:nn*no]

		copy(values[i*ns*no:(i+1)*ns*no], solver.Evaluate(&s, points))
		copy(moments[i*no:(i+1)*no], solver.Integrate(&s))
	}

	if config.Verbose {
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

func generate(problem *internal.Problem, target internal.Target) ([]float64, error) {
	config := &problem.Config.Assessment

	ni, _ := target.Dimensions()

	seed := int64(config.Seed)
	if seed < 0 {
		seed = time.Now().Unix()
	}

	ns := config.Samples
	if ns == 0 {
		return nil, errors.New("the number of samples should be positive")
	}

	return sequence.NewSobol(ni, seed).Next(ns), nil
}
