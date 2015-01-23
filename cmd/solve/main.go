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

	problem.Log(problem)
	problem.Log(target)

	var surrogate *adhier.Surrogate

	problem.Log("Constructing a surrogate...")
	problem.Log("Done in %v.", internal.Track(func() {
		surrogate = solver.Construct()
	}))

	problem.Log(surrogate)

	if f == nil {
		return nil
	}

	if err := f.Put("surrogate", *surrogate); err != nil {
		return err
	}

	return nil
}
