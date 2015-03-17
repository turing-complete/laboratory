package internal

import (
	"encoding/json"
	"os"

	"github.com/ready-steady/numeric/interpolation/adhier"
	"github.com/ready-steady/simulation/temperature/numeric"
)

// Config is a configuration of a problem.
type Config struct {
	Inherit string

	// A file containing the specification of a system (a platform and an
	// application) to analyze in the TGFF format.
	System string

	Probability   ProbabilityConfig   // Probability model
	Target        TargetConfig        // Quantity of interest
	Temperature   TemperatureConfig   // Temperature analysis
	Interpolation InterpolationConfig // Interpolation
	Assessment    AssessmentConfig    // Assessment

	// A flag for displaying diagnostic information.
	Verbose bool
}

// ProbabilityConfig is a configuration of the probability model.
type ProbabilityConfig struct {
	// The indices of the tasks whose execution times should be considered as
	// uncertain; if empty, the parameter is set to all tasks.
	TaskIndex []uint
	// The multiplier used to calculate the maximal delay of a task.
	MaxDelay float64 // ∈ [0, 1)
	// The marginal distributions of tasks’ delays.
	Marginal string
	// The strength of correlations between tasks.
	CorrLength float64 // > 0
	// The portion of the variance to be preserved when reducing the number of
	// stochastic dimensions.
	VarThreshold float64 // ∈ (0, 1]
}

// TargetConfig is a configuration of a quantity of interest.
type TargetConfig struct {
	// The name of the quantity of interest. The options are “end-to-end-delay,”
	// “total-energy,” and “temperature-profile.”
	Name string

	// The absolute error tolerance.
	Tolerance float64

	// The indices of the cores that should be considered; if empty, the
	// parameter is set to all cores.
	CoreIndex []uint
	// The time step of temperature profiles.
	TimeStep float64
	// The fraction of the application’s span that should be considered; if
	// empty, the parameter is set to the entire span [0, 1].
	TimeFraction []float64

	// A flag for displaying diagnostic information.
	Verbose bool
}

// TemperatureConfig is a configuration of the temperature analysis.
type TemperatureConfig struct {
	numeric.Config
}

// InterpolationConfig is a configuration of the interpolation algorithm.
type InterpolationConfig struct {
	// The quadrature rule to use, which is either “open” or “closed.”
	Rule string

	adhier.Config
}

// AssessmentConfig is a configuration of the assessment procedure.
type AssessmentConfig struct {
	// The seed for random number generation.
	Seed int
	// The number of samples to take.
	Samples uint
}

func NewConfig(path string) (Config, error) {
	paths := []string{path}
	for {
		config := Config{}
		if err := populate(&config, path); err != nil {
			return Config{}, err
		}

		if len(config.Inherit) > 0 {
			path = config.Inherit
			paths = append([]string{path}, paths...)
			continue
		}

		if len(paths) == 1 {
			return config, nil
		}

		break
	}

	config := Config{}
	for _, path := range paths {
		if err := populate(&config, path); err != nil {
			return Config{}, err
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
