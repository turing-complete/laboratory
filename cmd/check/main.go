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

	target, solver, err := problem.Setup()
	if err != nil {
		return err
	}

	problem.Log(problem)
	problem.Log(target)

	surrogate := new(adhier.Surrogate)
	if fi == nil {
		return errors.New("an input file is required")
	}
	if err := fi.Get("surrogate", surrogate); err != nil {
		return err
	}

	problem.Log(surrogate)

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

	var values, realValues []float64

	problem.Log("Evaluating the surrogate model...")
	problem.Log("Done in %v.", internal.Track(func() {
		values = solver.Evaluate(surrogate, points)
	}))

	problem.Log("Evaluating the original model...")
	problem.Log("Done in %v.", internal.Track(func() {
		realValues = solver.Compute(points)
	}))

	fmt.Printf("NRMSE: %e\n", metric.NRMSE(values, realValues))

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
