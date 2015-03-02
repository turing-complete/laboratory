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

	nt, ns := config.Assessment.Steps, config.Assessment.Samples
	if nt == 0 {
		nt = 1
	}

	no := uint(len(observations)) / (nt * ns)
	ni := uint(len(observationPoints)) / (nt * ns)

	if err := output.PutArray("observations", observations, no, ns, nt); err != nil {
		return err
	}
	if err := output.PutArray("observationPoints", observationPoints, ni, ns, nt); err != nil {
		return err
	}

	ni = uint(len(predictionPoints)) / (nt * ns)

	if err := output.PutArray("predictions", predictions, no, ns, nt); err != nil {
		return err
	}
	if err := output.PutArray("predictionPoints", predictionPoints, ni, ns, nt); err != nil {
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

	nt, ns := config.Steps, config.Samples
	if nt == 0 {
		nt = 1
	}
	if ns == 0 {
		return nil, errors.New("the number of samples is zero")
	}

	distribution := uniform.New(0, 1)

	ni, np := uint(target.Inputs()), uint(target.Pseudos())

	var fixed []float64

	if np > 0 {
		// If there are deterministic dimensions like time, we need to fix them
		// in order to generate comparable datasets. These dimensions are fixed
		// to randomly generated numbers, and this procedure is repeated
		// multiple times (specified by Steps) for a more comprehensive
		// assessment later on. The following line should be executed after the
		// seeding above and before the actual sampling below to ensure that it
		// chooses the same values each time this function is called.
		fixed = probability.Sample(distribution, nt*np)
	}

	samples := probability.Sample(distribution, nt*ns*ni)

	if np > 0 {
		for i := uint(0); i < nt; i++ {
			for j := uint(0); j < ns; j++ {
				for k := uint(0); k < np; k++ {
					samples[i*ns*ni+j*ni+k] = fixed[i*np+k]
				}
			}
		}
	}

	return samples, nil
}

func invoke(target internal.Target, points []float64) []float64 {
	nw := uint(runtime.GOMAXPROCS(0))
	ni, no := target.Inputs(), target.Outputs()
	np := uint(len(points)) / ni

	values := make([]float64, np*no)
	jobs := make(chan uint, np)
	group := sync.WaitGroup{}
	group.Add(int(np))

	for i := uint(0); i < nw; i++ {
		go func() {
			for j := range jobs {
				target.Evaluate(points[j*ni:(j+1)*ni], values[j*no:(j+1)*no], nil)
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
