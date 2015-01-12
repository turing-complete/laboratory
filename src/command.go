package main

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/ready-steady/format/mat"
	"github.com/ready-steady/numan/interp/adhier"
	"github.com/ready-steady/prob"
	"github.com/ready-steady/prob/uniform"
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

func show(p *problem, f *mat.File, _ *mat.File) error {
	fmt.Println(p)

	if f == nil {
		return nil
	}

	s := new(adhier.Surrogate)
	if err := f.Get("surrogate", s); err != nil {
		return err
	}

	fmt.Println(s)

	return nil
}

func solve(p *problem, _ *mat.File, f *mat.File) error {
	fmt.Println(p)

	var s *adhier.Surrogate
	track("Constructing a surrogate...", true, func() {
		s = p.solve()
	})

	fmt.Println(s)
	fmt.Println(p.cache)

	if f == nil {
		return nil
	}

	if err := f.Put("surrogate", *s); err != nil {
		return err
	}

	return nil
}

func check(p *problem, fi *mat.File, fo *mat.File) error {
	s := new(adhier.Surrogate)
	if fi == nil {
		return errors.New("an input file is required")
	}
	if err := fi.Get("surrogate", s); err != nil {
		return err
	}

	fmt.Println(p)
	fmt.Println(s)

	c := &p.config

	if c.Samples == 0 {
		return errors.New("the number of samples is zero")
	}

	rand.Seed(c.Seed)
	points := prob.Sample(uniform.New(0, 1), c.Samples*p.ic)

	var values, realValues []float64

	track("Evaluating the surrogate model...", true, func() {
		values = p.evaluate(s, points)
	})

	track("Evaluating the original model...", true, func() {
		realValues = p.compute(points)
	})

	fmt.Printf("NRMSE: %e\n", assess.NRMSE(values, realValues))

	if fo == nil {
		return nil
	}

	if err := fo.PutMatrix("points", points, p.ic, c.Samples); err != nil {
		return err
	}

	if err := fo.PutMatrix("values", values, p.oc, c.Samples); err != nil {
		return err
	}

	if err := fo.PutMatrix("realValues", realValues, p.oc, c.Samples); err != nil {
		return err
	}

	return nil
}
