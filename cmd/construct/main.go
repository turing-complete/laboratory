package main

import (
	"fmt"

	"github.com/ready-steady/format/mat"

	"../internal"
)

func main() {
	internal.Run(command)
}

func command(config internal.Config, _ *mat.File, output *mat.File) error {
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

	var solution *internal.Solution

	if config.Verbose {
		fmt.Println(problem)
		fmt.Println(target)
		fmt.Println("Constructing a surrogate...")
		solution = solver.Compute(target)
		fmt.Println(solution)
	} else {
		solution = solver.Compute(target)
	}

	if output == nil {
		return nil
	}

	if err := output.Put("solution", *solution); err != nil {
		return err
	}

	return nil
}
