package uncertainty

import (
	"errors"

	"github.com/ready-steady/statistics/correlation"
	"github.com/simulated-reality/laboratory/src/internal/config"
	"github.com/simulated-reality/laboratory/src/internal/system"

	icorrelation "github.com/simulated-reality/laboratory/src/internal/correlation"
)

type Uncertainty interface {
	Len() int
	Transform([]float64) []float64
}

func computeCorrelator(c *config.Uncertainty, s *system.System,
	taskIndex []uint) ([]float64, error) {

	if c.CorrLength < 0 {
		return nil, errors.New("the correlation length should be nonnegative")
	}
	if c.VarThreshold <= 0 {
		return nil, errors.New("the variance-reduction threshold should be positive")
	}

	C := icorrelation.Compute(s.Application, taskIndex, c.CorrLength)
	correlator, _, err := correlation.Decompose(C, uint(len(taskIndex)), c.VarThreshold)
	if err != nil {
		return nil, err
	}

	return correlator, nil
}
