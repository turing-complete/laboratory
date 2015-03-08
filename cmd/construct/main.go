package main

import (
	"fmt"

	"github.com/ready-steady/format/mat"
	"github.com/ready-steady/numeric/interpolation/adhier"

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

	var surrogate *adhier.Surrogate

	if config.Verbose {
		fmt.Println(problem)
		fmt.Println(target)
		fmt.Println("Constructing a surrogate...")

		surrogate = interpolator.Compute(target)
		target.Monitor(surrogate.Level, 0, surrogate.Nodes)

		fmt.Println(surrogate)
	} else {
		surrogate = interpolator.Compute(target)
	}

	if output == nil {
		return nil
	}

	if err := output.Put("surrogate", *surrogate); err != nil {
		return err
	}

	return nil
}
