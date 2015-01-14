package main

import (
	"github.com/ready-steady/numan/interp/adhier"
)

type solver interface {
	Construct() *adhier.Surrogate
	Compute(nodes []float64) []float64
	Evaluate(surrogate *adhier.Surrogate, points []float64) []float64
}

func newSolver(problem *problem, target target) solver {
	base := newBaseSolver(problem, target)

	if fc := base.ic - problem.zc; fc == 0 {
		return &directSolver{base}
	} else {
		return &cachedSolver{base, fc}
	}
}
