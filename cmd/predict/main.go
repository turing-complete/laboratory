package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/ready-steady/hdf5"
	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"

	"../internal"
)

func main() {
	internal.Run(command)
}

func command(config *internal.Config, input *hdf5.File, output *hdf5.File) error {
	if input == nil {
		return errors.New("an input file is required")
	}

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
	if err = input.Get("solution", solution); err != nil {
		return err
	}

	points, err := generate(problem, target)
	if err != nil {
		return err
	}

	if config.Verbose {
		fmt.Println("Sampling the surrogate model...")
		fmt.Println(problem)
		fmt.Println(target)
		fmt.Println(solution)
	}

	values := solver.Evaluate(solution, points)

	if config.Verbose {
		fmt.Println("Done.")
	}

	if output == nil {
		return nil
	}

	ns := config.Assessment.Samples
	no := uint(len(values)) / ns
	ni := uint(len(points)) / ns

	if err := output.Put("values", values, no, ns); err != nil {
		return err
	}
	if err := output.Put("points", points, ni, ns); err != nil {
		return err
	}
	if err := output.Put("solution", *solution); err != nil {
		return err
	}

	return nil
}

func generate(problem *internal.Problem, target internal.Target) ([]float64, error) {
	ni, _ := target.Dimensions()

	config := &problem.Config.Assessment

	if config.Seed > 0 {
		rand.Seed(int64(config.Seed))
	} else {
		rand.Seed(time.Now().Unix())
	}

	ns := config.Samples
	if ns == 0 {
		return nil, errors.New("the number of samples should be positive")
	}

	return probability.Sample(uniform.New(0, 1), ns*ni), nil
}
