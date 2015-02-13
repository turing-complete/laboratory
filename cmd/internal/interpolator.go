package internal

import (
	"errors"
	"strings"

	"github.com/ready-steady/numeric/basis/linhat"
	"github.com/ready-steady/numeric/grid/newcot"
	"github.com/ready-steady/numeric/interpolation/adhier"
)

func NewInterpolator(problem *Problem, target Target) (*adhier.Interpolator, error) {
	config := &problem.Config.Interpolation
	ic, oc := uint16(target.Inputs()), uint16(target.Outputs())

	var grid adhier.Grid
	var basis adhier.Basis

	switch strings.ToLower(config.Rule) {
	case "open":
		grid, basis = newcot.NewOpen(ic), linhat.NewOpen(ic)
	case "closed":
		grid, basis = newcot.NewClosed(ic), linhat.NewClosed(ic)
	default:
		return nil, errors.New("the interpolation rule is unknown")
	}

	// FIXME: Altering the problem’s data might not be a good idea.
	config.Inputs, config.Outputs = ic, oc

	return adhier.New(grid, basis, (*adhier.Config)(&config.Config))
}
