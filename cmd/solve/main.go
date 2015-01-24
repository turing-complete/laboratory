package main

import (
	"github.com/ready-steady/format/mat"
	"github.com/ready-steady/numeric/interpolation/adhier"

	"../internal"
)

func main() {
	internal.Run(command)
}

func command(_ *internal.Config, problem *internal.Problem,
	_ *mat.File, f *mat.File) error {

	target, solver, err := problem.Setup()
	if err != nil {
		return err
	}

	problem.Println(problem)
	problem.Println(target)

	var surrogate *adhier.Surrogate

	problem.Println("Constructing a surrogate...")
	problem.Printf("Done in %v.\n", internal.Track(func() {
		surrogate = solver.Construct()
	}))

	problem.Println(surrogate)

	if f == nil {
		return nil
	}

	if err := f.Put("surrogate", *surrogate); err != nil {
		return err
	}

	return nil
}
