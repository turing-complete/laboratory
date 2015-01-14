package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/ready-steady/format/mat"
	"github.com/ready-steady/numan/interp/adhier"
	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"
	"github.com/ready-steady/stats/assess"
)

func findCommand(name string) func(*problem, *mat.File, *mat.File) error {
	switch name {
	case "show":
		return show
	case "solve":
		return solve
	case "check":
		return check
	default:
		return nil
	}
}

func show(problem *problem, f *mat.File, _ *mat.File) error {
	fmt.Println(problem)

	if f == nil {
		return nil
	}

	surrogate := new(adhier.Surrogate)
	if err := f.Get("surrogate", surrogate); err != nil {
		return err
	}

	fmt.Println(surrogate)

	return nil
}

func solve(problem *problem, _ *mat.File, f *mat.File) error {
	target, solver, err := problem.setup()
	if err != nil {
		return err
	}

	fmt.Println(problem)
	fmt.Println(target)

	var surrogate *adhier.Surrogate
	track("Constructing a surrogate...", true, func() {
		surrogate = solver.Construct()
	})

	fmt.Println(surrogate)

	if f == nil {
		return nil
	}

	if err := f.Put("surrogate", *surrogate); err != nil {
		return err
	}

	return nil
}

func check(problem *problem, fi *mat.File, fo *mat.File) error {
	target, solver, err := problem.setup()
	if err != nil {
		return err
	}

	fmt.Println(problem)
	fmt.Println(target)

	surrogate := new(adhier.Surrogate)
	if fi == nil {
		return errors.New("an input file is required")
	}
	if err := fi.Get("surrogate", surrogate); err != nil {
		return err
	}

	fmt.Println(surrogate)

	sc := problem.config.Samples
	if sc == 0 {
		return errors.New("the number of samples is zero")
	}

	ic, oc := target.InputsOutputs()

	if problem.config.Seed > 0 {
		rand.Seed(problem.config.Seed)
	} else {
		rand.Seed(time.Now().Unix())
	}
	points := probability.Sample(uniform.New(0, 1), sc*ic)

	var values, realValues []float64

	track("Evaluating the surrogate model...", true, func() {
		values = solver.Evaluate(surrogate, points)
	})

	track("Evaluating the original model...", true, func() {
		realValues = solver.Compute(points)
	})

	fmt.Printf("NRMSE: %e\n", assess.NRMSE(values, realValues))

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
