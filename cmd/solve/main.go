package main

import (
	"github.com/ready-steady/format/mat"

	"../internal"
)

func main() {
	internal.Run(command)
}

func command(_ *internal.Config, problem *internal.Problem,
	_ *mat.File, f *mat.File) error {

	target, interpolator, err := internal.Setup(problem)
	if err != nil {
		return err
	}

	problem.Println(problem)
	problem.Println(target)

	problem.Println("Constructing a surrogate...")
	surrogate := interpolator.Compute(target.Evaluate)
	problem.Println(surrogate)

	if f == nil {
		return nil
	}

	if err := f.Put("surrogate", *surrogate); err != nil {
		return err
	}

	return nil
}
