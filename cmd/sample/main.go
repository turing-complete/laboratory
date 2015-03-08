package main

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/ready-steady/format/mat"
	"github.com/ready-steady/numeric/interpolation/adhier"
	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"

	"../internal"
)

var startTime = time.Now().Unix()

func main() {
	internal.Run(command)
}

func command(config internal.Config, input *mat.File, output *mat.File) error {
	if input == nil {
		return errors.New("an input file is required")
	}

	if config.Verbose {
		fmt.Println("Processing the original model...")
	}
	observations, observationPoints, err := observe(config)
	if err != nil {
		return err
	}
	if config.Verbose {
		fmt.Println("Done.")
	}

	if config.Verbose {
		fmt.Println("Processing the surrogate model...")
	}
	predictions, predictionPoints, surrogate, err := predict(config, input)
	if err != nil {
		return err
	}
	if config.Verbose {
		fmt.Println("Done.")
	}

	if output == nil {
		return nil
	}

	ns := config.Assessment.Samples
	no := uint(len(observations)) / ns
	ni := uint(len(observationPoints)) / ns

	if err := output.PutArray("observations", observations, no, ns); err != nil {
		return err
	}
	if err := output.PutArray("observationPoints", observationPoints, ni, ns); err != nil {
		return err
	}

	ni = uint(len(predictionPoints)) / ns

	if err := output.PutArray("predictions", predictions, no, ns); err != nil {
		return err
	}
	if err := output.PutArray("predictionPoints", predictionPoints, ni, ns); err != nil {
		return err
	}

	if err := output.Put("surrogate", *surrogate); err != nil {
		return err
	}

	return nil
}

func observe(config internal.Config) ([]float64, []float64, error) {
	config.Probability.VarThreshold = 42

	problem, err := internal.NewProblem(config)
	if err != nil {
		return nil, nil, err
	}

	target, err := internal.NewTarget(problem)
	if err != nil {
		return nil, nil, err
	}

	points, err := generate(problem, target)
	if err != nil {
		return nil, nil, err
	}

	if config.Verbose {
		fmt.Println(problem)
		fmt.Println(target)
	}

	return invoke(target, points), points, nil
}

func predict(config internal.Config, input *mat.File) (
	[]float64, []float64, *adhier.Surrogate, error) {

	problem, err := internal.NewProblem(config)
	if err != nil {
		return nil, nil, nil, err
	}

	target, err := internal.NewTarget(problem)
	if err != nil {
		return nil, nil, nil, err
	}

	interpolator, err := internal.NewInterpolator(problem, target)
	if err != nil {
		return nil, nil, nil, err
	}

	surrogate := new(adhier.Surrogate)
	if err = input.Get("surrogate", surrogate); err != nil {
		return nil, nil, nil, err
	}

	if config.Verbose {
		fmt.Println(problem)
		fmt.Println(target)
		fmt.Println(surrogate)
	}

	points, err := generate(problem, target)
	if err != nil {
		return nil, nil, nil, err
	}

	return interpolator.Evaluate(surrogate, points), points, surrogate, nil
}

func generate(problem *internal.Problem, target internal.Target) ([]float64, error) {
	config := &problem.Config.Assessment

	if config.Seed > 0 {
		rand.Seed(int64(config.Seed))
	} else {
		rand.Seed(startTime)
	}

	ns := config.Samples
	if ns == 0 {
		return nil, errors.New("the number of samples should be positive")
	}

	ni, _ := target.Dimensions()

	return probability.Sample(uniform.New(0, 1), ns*ni), nil
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
