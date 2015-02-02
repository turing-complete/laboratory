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
	problem.Printf("%10s %15s %15s\n", "Level", "Nodes", "Evaluations")

	report := func(level uint8, nodes uint32) {
		problem.Printf("%10d %15d %15d\n", level, nodes, target.Evaluations())
	}
	progress := func(level uint8, activeNodes uint32, totalNodes uint32) {
		if level > 0 {
			report(level-1, totalNodes-activeNodes)
		}
	}

	surrogate := interpolator.Compute(target.Evaluate, progress)
	report(surrogate.Level, surrogate.Nodes)
	problem.Println(surrogate)

	if f == nil {
		return nil
	}

	if err := f.Put("surrogate", *surrogate); err != nil {
		return err
	}

	return nil
}
