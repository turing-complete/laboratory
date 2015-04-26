package internal

import (
	"encoding/json"
	"os"

	"github.com/ready-steady/adapt"
	"github.com/ready-steady/simulation/temperature/analytic"
)

// Config is a configuration of the problem.
type Config struct {
	Inherit string

	System        SystemConfig        // Platform and application
	Probability   ProbabilityConfig   // Probability model
	Target        TargetConfig        // Quantity of interest
	Interpolation InterpolationConfig // Interpolation
	Assessment    AssessmentConfig    // Assessment

	// A flag indicating that diagnostic information should be displayed.
	Verbose bool
}

// SystemConfig is a configuration of the system.
type SystemConfig struct {
	// A TGFF file describing the platform and application to analyze.
	Specification string

	analytic.Config
}

// ProbabilityConfig is a configuration of the probability model.
type ProbabilityConfig struct {
	// The tasks whose execution times should be considered as uncertain.
	TaskIndex string // ⊂ {0, ..., #tasks-1}

	// The seed for initializing the tasks’ execution modes.
	Seed int64
	// The number of modes per task.
	Modes uint // > 0
	// The minimal scaling factor of a mode.
	MinScale float64 // > 0
	// The maximal scaling factor of a mode.
	MaxScale float64 // > 0
	// The relative length of transition from one mode to another.
	Transition float64 // ∈ (0, 0.5]

	// The strength of correlations between tasks.
	CorrLength float64 // > 0
	// The portion of the variance to be preserved when reducing the number of
	// stochastic dimensions.
	VarThreshold float64 // ∈ (0, 1]
}

// TargetConfig is a configuration of the quantity of interest.
type TargetConfig struct {
	// The name of the quantity. The options are “end-to-end-delay,”
	// “total-energy,” and “temperature-profile.”
	Name string

	// The error tolerance.
	Tolerance float64
	// The patter that is replicated onto a surplus in order to identify the
	// elements that should be used for the error estimation.
	Stencil []bool

	// The cores that should be considered.
	CoreIndex string // ⊂ {0, ..., #cores-1}
	// The time moments that should be considered. The elements are assumed to
	// be normalized by the application’s span.
	TimeIndex string // ⊂ [0, 1]

	// A flag indicating that diagnostic information should be displayed.
	Verbose bool
}

// InterpolationConfig is a configuration of the interpolation algorithm.
type InterpolationConfig struct {
	// The quadrature rule to use, which is either “closed” or “open.”
	Rule string

	adapt.Config
}

// AssessmentConfig is a configuration of the assessment procedure.
type AssessmentConfig struct {
	// The seed for generating samples.
	Seed int64
	// The number of samples to draw.
	Samples uint
	// The number of bins used for comparing distributions.
	Bins uint
}

func NewConfig(path string) (*Config, error) {
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
