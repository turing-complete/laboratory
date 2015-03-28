package internal

import (
	"github.com/ready-steady/ode/dopri"
	"github.com/ready-steady/simulation/temperature/numeric"
)

func newTemperature(config *TemperatureConfig) (*numeric.Temperature, error) {
	integrator, err := dopri.New(&dopri.Config{
		MaxStep:  0,
		TryStep:  0,
		AbsError: 1e-3,
		RelError: 1e-3,
	})
	if err != nil {
		return nil, err
	}

	return numeric.New(&config.Config, integrator), nil
}
