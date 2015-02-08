package main

import (
	"errors"
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

	observed, err := observe(config)
	if err != nil {
		return err
	}

	predicted, surrogate, err := predict(config, input)
	if err != nil {
		return err
	}

	if output == nil {
		return nil
	}

	if err := output.Put("surrogate", *surrogate); err != nil {
		return err
	}

	oc := surrogate.Outputs
	sc := uint32(len(predicted)) / oc

	if err := output.PutMatrix("predicted", predicted, oc, sc); err != nil {
		return err
	}
	if err := output.PutMatrix("observed", observed, oc, sc); err != nil {
		return err
	}

	return nil
}

func observe(config internal.Config) ([]float64, error) {
	config.ProbModel.VarThreshold = 42

	problem, err := internal.NewProblem(config)
	if err != nil {
		return nil, err
	}

	target, err := internal.NewTarget(problem)
	if err != nil {
		return nil, err
	}

	points, err := generate(problem, target)
	if err != nil {
		return nil, err
	}

	problem.Println("Processing the original model...")

	problem.Println(problem)
	problem.Println(target)

	return invoke(target, points), nil
}

func predict(config internal.Config, input *mat.File) ([]float64, *adhier.Surrogate, error) {
	problem, err := internal.NewProblem(config)
	if err != nil {
		return nil, nil, err
	}

	target, err := internal.NewTarget(problem)
	if err != nil {
		return nil, nil, err
	}

	interpolator, err := internal.NewInterpolator(problem, target)
	if err != nil {
		return nil, nil, err
	}

	surrogate := new(adhier.Surrogate)
	if err = input.Get("surrogate", surrogate); err != nil {
		return nil, nil, err
	}

	problem.Println("Processing the surrogate model...")

	problem.Println(problem)
	problem.Println(target)
	problem.Println(surrogate)

	points, err := generate(problem, target)
	if err != nil {
		return nil, nil, err
	}

	return interpolator.Evaluate(surrogate, points), surrogate, nil
}

func generate(problem *internal.Problem, target internal.Target) ([]float64, error) {
	config := &problem.Config.Assessment

	ic := target.Inputs()

	sc := config.Samples
	if sc == 0 {
		return nil, errors.New("the number of samples is zero")
	}

	if config.Seed > 0 {
		rand.Seed(config.Seed)
	} else {
		rand.Seed(startTime)
	}

	return probability.Sample(uniform.New(0, 1), sc*ic), nil
}

func invoke(target internal.Target, points []float64) []float64 {
	wc := uint32(runtime.GOMAXPROCS(0))
	ic, oc := target.Inputs(), target.Outputs()
	pc := uint32(len(points)) / ic

	values := make([]float64, pc*oc)
	jobs := make(chan uint32, pc)
	group := sync.WaitGroup{}
	group.Add(int(pc))

	for i := uint32(0); i < wc; i++ {
		go func() {
			for j := range jobs {
				target.Evaluate(points[j*ic:(j+1)*ic], values[j*oc:(j+1)*oc], nil)
				group.Done()
			}
		}()
	}

	for i := uint32(0); i < pc; i++ {
		jobs <- i
	}

	group.Wait()
	close(jobs)

	return values
}
