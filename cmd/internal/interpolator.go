package internal

import (
	"errors"
	"strings"

	"github.com/ready-steady/numeric/basis/linhat"
	"github.com/ready-steady/numeric/grid/newcot"
	"github.com/ready-steady/numeric/interpolation/adhier"
)

func NewInterpolator(problem *Problem, target Target) (*adhier.Interpolator, error) {
	ni, no := target.Inputs(), target.Outputs()

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

	config := (adhier.Config)(problem.Config.Interpolation.Config)
	config.Inputs, config.Outputs = ni, no

	return adhier.New(grid, basis, &config)
}
