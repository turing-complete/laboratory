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

	target, interpolator, err := internal.Setup(problem)
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

	problem.Println("Evaluating the surrogate model...")
	values := interpolator.Evaluate(surrogate, points)

	problem.Println("Evaluating the original model...")
	realValues := make([]float64, sc*oc)
	done := make(chan bool, sc)
	for i := uint32(0); i < sc; i++ {
		go func(point, value []float64) {
			target.Evaluate(point, value, nil)
			done <- true
		}(points[i*ic:(i+1)*ic], realValues[i*oc:(i+1)*oc])
	}
	for i := uint32(0); i < sc; i++ {
		<-done
	}

	fmt.Printf("NRMSE: %.2e\n", metric.NRMSE(values, realValues))

	if fo == nil {
		return nil
	}

	if err := fo.PutMatrix("points", points, ic, sc); err != nil {
		return err
	}
	if err := fo.PutMatrix("values", values, oc, sc); err != nil {
		return err
	}
	if err := fo.PutMatrix("realValues", realValues, oc, sc); err != nil {
		return err
	}

	return nil
}
