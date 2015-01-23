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

func command(_ *internal.Config, problem *internal.Problem,
	f *mat.File, _ *mat.File) error {

	fmt.Println(problem)

	if f == nil {
		return nil
	}

	surrogate := new(adhier.Surrogate)
	if err := f.Get("surrogate", surrogate); err != nil {
		return err
	}

	fmt.Println(surrogate)

	return nil
}
