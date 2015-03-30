package main

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/ready-steady/hdf5"

	"../internal"
)

func main() {
	internal.Run(command)
}

func command(config *internal.Config, _ *hdf5.File, output *hdf5.File) error {
	config.Probability.VarThreshold = 42

	problem, err := internal.NewProblem(config)
	if err != nil {
		return err
	}

	target, err := internal.NewTarget(problem)
	if err != nil {
		return err
	}

	points, err := generate(problem, target)
	if err != nil {
		return err
	}

	if config.Verbose {
		fmt.Println("Sampling the original model...")
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

	ns := config.Assessment.Samples
	no := uint(len(values)) / ns
	ni := uint(len(points)) / ns

	if err := output.Put("values", values, no, ns); err != nil {
		return err
	}
	if err := output.Put("points", points, ni, ns); err != nil {
		return err
	}

	return nil
}

func generate(problem *internal.Problem, target internal.Target) ([]float64, error) {
	config := &problem.Config.Assessment

	if config.Seed > 0 {
		rand.Seed(int64(config.Seed))
	} else {
		rand.Seed(time.Now().Unix())
	}

	ns := config.Samples
	if ns == 0 {
		return nil, errors.New("the number of samples should be positive")
	}

	return target.Generate(ns), nil
}
