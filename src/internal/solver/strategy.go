package solver

import (
	"github.com/turing-complete/laboratory/src/internal/config"

	algorithm "github.com/ready-steady/adapt/algorithm/local"
)

type strategy struct {
	algorithm.Strategy
}

func newStrategy(ni, no uint, config *config.Solver, grid algorithm.Grid) *strategy {
	return &strategy{*algorithm.NewStrategy(ni, no, config.MinLevel,
		config.MaxLevel, config.LocalError, grid)}
}
