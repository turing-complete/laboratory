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

	tc, sc := config.Assessment.Steps, config.Assessment.Samples
	if tc == 0 {
		tc = 1
	}

	oc := uint(len(observations)) / (tc * sc)
	ic := uint(len(observationPoints)) / (tc * sc)

	if err := output.PutArray("observations", observations, oc, sc, tc); err != nil {
		return err
	}
	if err := output.PutArray("observationPoints", observationPoints, ic, sc, tc); err != nil {
		return err
	}

	ic = uint(len(predictionPoints)) / (tc * sc)

	if err := output.PutArray("predictions", predictions, oc, sc, tc); err != nil {
		return err
	}
	if err := output.PutArray("predictionPoints", predictionPoints, ic, sc, tc); err != nil {
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

	tc, sc := config.Steps, config.Samples
	if tc == 0 {
		tc = 1
	}
	if sc == 0 {
		return nil, errors.New("the number of samples is zero")
	}

	distribution := uniform.New(0, 1)

	ic, pc := uint(target.Inputs()), uint(target.Pseudos())

	var fixed []float64

	if pc > 0 {
		// If there are deterministic dimensions like time, we need to fix them
		// in order to generate comparable datasets. These dimensions are fixed
		// to randomly generated numbers, and this procedure is repeated
		// multiple times (specified by Steps) for a more comprehensive
		// assessment later on. The following line should be executed after the
		// seeding above and before the actual sampling below to ensure that it
		// chooses the same values each time this function is called.
		fixed = probability.Sample(distribution, tc*pc)
	}

	samples := probability.Sample(distribution, tc*sc*ic)

	if pc > 0 {
		for i := uint(0); i < tc; i++ {
			for j := uint(0); j < sc; j++ {
				for k := uint(0); k < pc; k++ {
					samples[i*sc*ic+j*ic+k] = fixed[i*pc+k]
				}
			}
		}
	}

	return samples, nil
}

func invoke(target internal.Target, points []float64) []float64 {
	wc := uint(runtime.GOMAXPROCS(0))
	ic, oc := target.Inputs(), target.Outputs()
	pc := uint(len(points)) / ic

	values := make([]float64, pc*oc)
	jobs := make(chan uint, pc)
	group := sync.WaitGroup{}
	group.Add(int(pc))

	for i := uint(0); i < wc; i++ {
		go func() {
			for j := range jobs {
				target.Evaluate(points[j*ic:(j+1)*ic], values[j*oc:(j+1)*oc], nil)
				group.Done()
			}
		}()
	}

	for i := uint(0); i < pc; i++ {
		jobs <- i
	}

	group.Wait()
	close(jobs)

	return values
}
