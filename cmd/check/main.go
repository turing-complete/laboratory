package main

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/ready-steady/format/mat"
	"github.com/ready-steady/numeric/interpolation/adhier"
	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"
	"github.com/ready-steady/statistics/test"

	"../internal"
)

func main() {
	internal.Run(command)
}

func command(config internal.Config, input *mat.File, output *mat.File) error {
	approximations, surrogate, err := sampleSurrogate(config, input)
	if err != nil {
		return err
	}

	values, err := sampleOriginal(config)
	if err != nil {
		return err
	}

	α := config.Assessment.Alpha
	if α <= 0 || α > 1 {
		α = 0.05
	}

	rejected, p, Δ := test.KolmogorovSmirnov(approximations, values, α)

	fmt.Printf("Inputs: %d, outputs: %d, level: %2d, nodes: %6d, rejected: %5v"+
		" (α %.3f, p %.3e, Δ %.3e)\n", surrogate.Inputs, surrogate.Outputs,
		surrogate.Level, surrogate.Nodes, rejected, α, p, Δ)

	if output == nil {
		return nil
	}

	oc := surrogate.Outputs
	sc := uint32(len(approximations)) / oc

	if err := output.PutMatrix("approximations", approximations, oc, sc); err != nil {
		return err
	}
	if err := output.PutMatrix("values", values, oc, sc); err != nil {
		return err
	}

	return nil
}

func sampleSurrogate(config internal.Config, input *mat.File) ([]float64, *adhier.Surrogate, error) {
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
	if input == nil {
		return nil, nil, errors.New("an input file is required")
	}
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

func sampleOriginal(config internal.Config) ([]float64, error) {
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

func generate(problem *internal.Problem, target internal.Target) ([]float64, error) {
	config := &problem.Config.Assessment

	sc := config.Samples
	if sc == 0 {
		return nil, errors.New("the number of samples is zero")
	}

	if config.Seed > 0 {
		rand.Seed(config.Seed)
	} else {
		rand.Seed(time.Now().Unix())
	}

	ic, _ := target.InputsOutputs()

	return probability.Sample(uniform.New(0, 1), sc*ic), nil
}

func invoke(target internal.Target, points []float64) []float64 {
	wc := uint32(runtime.GOMAXPROCS(0))
	ic, oc := target.InputsOutputs()
	pc := uint32(len(points)) / ic

	values := make([]float64, pc*oc)
	jobs := make(chan uint32, pc)
	done := make(chan bool, pc)

	for i := uint32(0); i < wc; i++ {
		go func() {
			for j := range jobs {
				target.Evaluate(points[j*ic:(j+1)*ic], values[j*oc:(j+1)*oc], nil)
				done <- true
			}
		}()
	}

	for i := uint32(0); i < pc; i++ {
		jobs <- i
	}
	for i := uint32(0); i < pc; i++ {
		<-done
	}

	close(jobs)

	return values
}
