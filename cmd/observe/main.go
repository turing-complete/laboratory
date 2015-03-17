package main

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/ready-steady/format/mat"

	"../internal"
)

func main() {
	internal.Run(command)
}

func command(config internal.Config, input *mat.File, output *mat.File) error {
	if input == nil {
		return errors.New("an input file is required")
	}

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

	values := invoke(target, points)

	if config.Verbose {
		fmt.Println("Done.")
	}

	if output == nil {
		return nil
	}

	ns := config.Assessment.Samples
	no := uint(len(values)) / ns
	ni := uint(len(points)) / ns

	if err := output.PutArray("values", values, no, ns); err != nil {
		return err
	}
	if err := output.PutArray("points", points, ni, ns); err != nil {
		return err
	}

	return nil
}

func invoke(target internal.Target, points []float64) []float64 {
	nw := uint(runtime.GOMAXPROCS(0))
	ni, no := target.Dimensions()
	np := uint(len(points)) / ni

	values := make([]float64, np*no)
	jobs := make(chan uint, np)
	group := sync.WaitGroup{}
	group.Add(int(np))

	for i := uint(0); i < nw; i++ {
		go func() {
			for j := range jobs {
				target.Compute(points[j*ni:(j+1)*ni], values[j*no:(j+1)*no])
				group.Done()
			}
		}()
	}

	for i := uint(0); i < np; i++ {
		jobs <- i
	}

	group.Wait()
	close(jobs)

	return values
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
