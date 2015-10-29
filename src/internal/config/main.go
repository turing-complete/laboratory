package config

import (
	"encoding/json"
	"os"

	"github.com/ready-steady/adapt"
	"github.com/turing-complete/temperature/analytic"
)

// Config is a configuration of the problem.
type Config struct {
	Inherit string

	System        System        // Platform and application
	Uncertainty   Uncertainty   // Probability model
	Target        Target        // Quantity of interest
	Interpolation Interpolation // Interpolation
	Assessment    Assessment    // Assessment

	// A flag to display diagnostic information.
	Verbose bool
}

// System is a configuration of the system.
type System struct {
	// A TGFF file describing the platform and application to analyze.
	Specification string

	analytic.Config
}

// Uncertainty is a configuration of the probability model.
type Uncertainty struct {
	// The tasks whose execution times should be considered as uncertain.
	TaskIndex string // ⊂ {0, ..., #tasks-1}

	// The seed for initializing the tasks’ execution modes.
	Seed int64
	// The number of modes per task.
	Modes uint // > 0
	// The minimal relative offset of a mode.
	MinOffset float64 // > -1
	// The maximal relative offset of a mode.
	MaxOffset float64 // > -1
	// The relative length of transition from one mode to another.
	Transition float64 // ∈ (0, 0.5]

	// The strength of correlations between tasks.
	CorrLength float64 // > 0
	// The portion of the variance to be preserved when reducing the number of
	// stochastic dimensions.
	VarThreshold float64 // ∈ (0, 1]
}

// Target is a configuration of the quantity of interest.
type Target struct {
	// The name of the quantity. The options are “end-to-end-delay,”
	// “total-energy,” and “temperature-profile.”
	Name string

	// The weights for output dimensions.
	Importance []float64
	// The rejection threshold for output dimensions.
	Rejection []float64
	// The refinement threshold for output dimensions.
	Refinement []float64

	// The cores that should be considered.
	CoreIndex string // ⊂ {0, ..., #cores-1}
	// The time moments that should be considered. The elements are assumed to
	// be normalized by the application’s span.
	TimeIndex string // ⊂ [0, 1]

	// A flag to display diagnostic information.
	Verbose bool
}

// Interpolation is a configuration of the interpolation algorithm.
type Interpolation struct {
	// The quadrature rule to use, which is either “closed” or “open.”
	Rule string

	adapt.Config
}

// Assessment is a configuration of the assessment procedure.
type Assessment struct {
	// A flag to use the analytically computed moments.
	Analytic []bool
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
