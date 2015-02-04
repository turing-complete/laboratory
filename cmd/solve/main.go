package main

import (
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

	interpolator, err := internal.NewInterpolator(problem, target)
	if err != nil {
		return err
	}

	problem.Println(problem)
	problem.Println(target)

	problem.Println("Constructing a surrogate...")
	surrogate := interpolator.Compute(target.Evaluate, target.Progress)
	target.Progress(surrogate.Level, 0, surrogate.Nodes)
	problem.Println(surrogate)

	if output == nil {
		return nil
	}

	if err := output.Put("surrogate", *surrogate); err != nil {
		return err
	}

	return nil
}
