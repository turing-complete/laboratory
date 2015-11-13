package config

import (
	"encoding/json"
	"os"

	"github.com/ready-steady/adapt"

	temperature "github.com/turing-complete/temperature/analytic"
)

// Config is a configuration of the problem.
type Config struct {
	Inherit string

	System     System     // Platform and application
	Target     Target     // Quantity of interest
	Solver     Solver     // Interpolation algorithm
	Assessment Assessment // Assessment

	// A flag to display diagnostic information.
	Verbose bool
}

// System is a configuration of the system.
type System struct {
	// A TGFF file describing the platform and application to analyze.
	Specification string

	temperature.Config
}

// Target is a configuration of the quantity of interest.
type Target struct {
	// The name of the quantity. The options are “end-to-end-delay,”
	// “total-energy,” and “maximal-temperature.”
	Name string

	/// The probability model.
	Uncertainty Uncertainty

	// The weights for output dimensions.
	Importance []float64
	// The rejection threshold for output dimensions.
	Rejection []float64
	// The refinement threshold for output dimensions.
	Refinement []float64
}

// Uncertainty is a configuration of the probability model.
type Uncertainty struct {
	// The tasks whose execution times should be considered as uncertain.
	Tasks string // ⊂ {0, ..., #tasks-1}

	// The marginal distributions of tasks’ delays.
	Distribution string
	// The multiplier used to calculate the deviation of a parameter.
	Deviation float64 // ≥ 0

	// The strength of correlations between tasks.
	CorrLength float64 // > 0
	// The portion of the variance to be preserved.
	Reduction float64 // ∈ (0, 1]
}

// Solver is a configuration of the interpolation algorithm.
type Solver struct {
	// The quadrature rule to use, which is either “closed” or “open.”
	Rule string

	adapt.Config
}

// Assessment is a configuration of the assessment procedure.
type Assessment struct {
	// The seed for generating samples.
	Seed int64
	// The number of samples to draw.
	Samples uint
}

func New(path string) (*Config, error) {
	paths := []string{path}
	for {
		config := Config{}
		if err := populate(&config, path); err != nil {
			return nil, err
		}

		if len(config.Inherit) > 0 {
			path = config.Inherit
			paths = append([]string{path}, paths...)
			continue
		}

		break
	}

	config := &Config{}
	for _, path := range paths {
		if err := populate(config, path); err != nil {
			return nil, err
		}
	}

	return config, nil
}

func populate(config *Config, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	return decoder.Decode(config)
}
