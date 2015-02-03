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
	ic, oc := target.InputsOutputs()

	var grid adhier.Grid
	var basis adhier.Basis

	switch strings.ToLower(config.Rule) {
	case "open":
		grid = newcot.NewOpen(uint16(ic))
		basis = linhat.NewOpen(uint16(ic), uint16(oc))
	case "closed":
		grid = newcot.NewClosed(uint16(ic))
		basis = linhat.NewClosed(uint16(ic), uint16(oc))
	default:
		return nil, errors.New("the interpolation rule is unknown")
	}

	return adhier.New(grid, basis, (*adhier.Config)(&config.Config))
}
