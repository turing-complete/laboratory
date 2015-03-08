package internal

import (
	"errors"
	"strings"

	"github.com/ready-steady/numeric/basis/linhat"
	"github.com/ready-steady/numeric/grid/newcot"
	"github.com/ready-steady/numeric/interpolation/adhier"
)

func NewInterpolator(problem *Problem, target Target) (*adhier.Interpolator, error) {
	ni, _ := target.Dimensions()

	var grid adhier.Grid
	var basis adhier.Basis

	switch strings.ToLower(problem.Config.Interpolation.Rule) {
	case "open":
		grid, basis = newcot.NewOpen(ni), linhat.NewOpen(ni)
	case "closed":
		grid, basis = newcot.NewClosed(ni), linhat.NewClosed(ni)
	default:
		return nil, errors.New("the interpolation rule is unknown")
	}

	return adhier.New(grid, basis, (*adhier.Config)(&problem.Config.Interpolation.Config)), nil
}
