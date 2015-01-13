package main

import (
	"testing"
)

func BenchmarkSolverConstruct(b *testing.B) {
	config, _ := loadConfig("fixtures/002_020.json")

	for i := 0; i < b.N; i++ {
		problem, _ := newProblem(config)
		target, _ := newTarget(problem)
		solver := newSolver(problem, target)
		solver.construct()
	}
}
