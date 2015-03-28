package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/ready-steady/hdf5"

	"../internal"
)

var numberOfPoints = flag.Uint("number", 100, "the number of points per dimension")

func main() {
	internal.Run(command)
}

func command(config internal.Config, _ *hdf5.File, output *hdf5.File) error {
	problem, err := internal.NewProblem(config)
	if err != nil {
		return err
	}

	target, err := internal.NewTarget(problem)
	if err != nil {
		return err
	}

	points := generate(target)

	if config.Verbose {
		fmt.Println("Sweeping the reduced model...")
		fmt.Println(problem)
		fmt.Println(target)
	}

	values := internal.Invoke(target, points, uint(runtime.GOMAXPROCS(0)))

	if config.Verbose {
		fmt.Println("Done.")
	}

	if output == nil {
		return nil
	}

	ni, no := target.Dimensions()
	np := uint(len(points)) / ni

	if err := output.Put("values", values, no, np); err != nil {
		return err
	}
	if err := output.Put("points", points, ni, np); err != nil {
		return err
	}

	return nil
}

func generate(target internal.Target) []float64 {
	np := *numberOfPoints
	ni, _ := target.Dimensions()

	points := make([]float64, np*ni)

	for i := uint(0); i < np; i++ {
		value := float64(i) * 1.0 / float64(np-1)
		for j := uint(0); j < ni; j++ {
			points[i*ni+j] = value
		}
	}

	return points
}
