package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/ready-steady/format/mat"
	"github.com/ready-steady/numeric/interpolation/adhier"
	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"
	"github.com/ready-steady/statistics/metric"

	"../internal"
)

func main() {
	internal.Run(command)
}

func command(config *internal.Config, problem *internal.Problem,
	fi *mat.File, fo *mat.File) error {

	target, err := internal.SetupTarget(problem)
	if err != nil {
		return err
	}

	interpolator, err := internal.SetupInterpolator(problem, target)
	if err != nil {
		return err
	}

	problem.Println(problem)
	problem.Println(target)

	surrogate := new(adhier.Surrogate)
	if fi == nil {
		return errors.New("an input file is required")
	}
	if err := fi.Get("surrogate", surrogate); err != nil {
		return err
	}

	problem.Println(surrogate)

	sc := config.Samples
	if sc == 0 {
		return errors.New("the number of samples is zero")
	}

	ic, oc := target.InputsOutputs()

	if config.Seed > 0 {
		rand.Seed(config.Seed)
	} else {
		rand.Seed(time.Now().Unix())
	}
	points := probability.Sample(uniform.New(0, 1), sc*ic)

	problem.Println("Evaluating the original model...")
	values := invoke(target, points)

	problem.Println("Evaluating the surrogate model...")
	approximations := interpolator.Evaluate(surrogate, points)

	fmt.Printf("NRMSE: %.2e\n", metric.NRMSE(approximations, values))

	if fo == nil {
		return nil
	}

	if err := fo.PutMatrix("points", points, ic, sc); err != nil {
		return err
	}
	if err := fo.PutMatrix("values", values, oc, sc); err != nil {
		return err
	}
	if err := fo.PutMatrix("approximations", approximations, oc, sc); err != nil {
		return err
	}

	return nil
}

func invoke(target internal.Target, points []float64) []float64 {
	ic, oc := target.InputsOutputs()
	pc := uint32(len(points)) / ic

	values := make([]float64, pc*oc)
	done := make(chan bool, pc)

	for i := uint32(0); i < pc; i++ {
		go func(point, value []float64) {
			target.Evaluate(point, value, nil)
			done <- true
		}(points[i*ic:(i+1)*ic], values[i*oc:(i+1)*oc])
	}
	for i := uint32(0); i < pc; i++ {
		<-done
	}

	return values
}
