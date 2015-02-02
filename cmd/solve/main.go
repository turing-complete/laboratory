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
	problem.Printf("%5s %10s %10s\n", "Level", "Active", "Total")
	progress := func(level uint8, activeNodes uint32, totalNodes uint32) {
		problem.Printf("%5d %10d %10d\n", level, activeNodes, totalNodes)
	}

	surrogate := interpolator.Compute(target.Evaluate, progress)
	problem.Println(surrogate)

	if f == nil {
		return nil
	}

	if err := f.Put("surrogate", *surrogate); err != nil {
		return err
	}

	return nil
}
