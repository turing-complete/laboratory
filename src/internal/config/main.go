package config

import (
	"encoding/json"
	"os"

	temperature "github.com/turing-complete/temperature/analytic"
)

// Config is a configuration of a problem.
type Config struct {
	Inherit string

	// The system
	System System
	// The quantity of interest
	Quantity Quantity
	// The probability model
	Uncertainty Uncertainty
	// The approximation algorithm
	Solution Solution
	// The assessment procedure
	Assessment Assessment
	// The flag for displaying diagnostic information.
	Verbose bool
}

// System is a configuration of the system.
type System struct {
	// The TGFF file describing the platform and application.
	Specification string

	temperature.Config
}

// Quantity is a configuration of the quantity of interest.
type Quantity struct {
	// The name of the quantity. The options are “end-to-end-delay,”
	// “total-energy,” and “maximum-temperature.”
	Name string
}

// Uncertainty is a configuration of the probability model.
type Uncertainty struct {
	// The tasks whose execution times should be considered as uncertain.
	Tasks string // ⊂ {0, …, #tasks-1}
	// The marginal distributions of tasks’ delays.
	Distribution string
	// The multiplier used to calculate the range of deviation.
	Deviation float64 // ≥ 0
	// The strength of correlations between tasks.
	Correlation float64 // > 0
	// The portion of the variance to be preserved.
	Variance float64 // ∈ (0, 1]
}

// Solution is a configuration of the approximation algorithm.
type Solution struct {
	// The flag for interpolating with the probability distribution of the
	// uncertain parameters embedded into the surrogate.
	Aleatory bool
	// The quadrature rule, which is either “closed” or “open.”
	Rule string
	// The total order of polynomials.
	Power uint
	// The minimum level of approximation.
	MinLevel uint
	// The maximum level of approximation.
	MaxLevel uint
	// The maximum number of evaluations.
	MaxEvaluations uint
	// The tolerance of the absolute error.
	AbsoluteError float64
	// The tolerance of the relative error.
	RelativeError float64
	// The tolerance of the score error.
	ScoreError float64
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
	return json.NewDecoder(file).Decode(config)
}
