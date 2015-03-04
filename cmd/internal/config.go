package internal

import (
	"encoding/json"
	"os"

	"github.com/ready-steady/numeric/interpolation/adhier"
	"github.com/ready-steady/simulation/temperature/numeric"
)

type Config struct {
	// The TGFF file of the system to analyze.
	TGFF string

	// The probability model.
	Probability struct {
		// The indices of the tasks whose execution times should be considered
		// as uncertain; if empty, the parameter is set to all tasks.
		Tasks []uint
		// The multiplier used to calculate the maximal delay of a task.
		MaxDelay float64 // ∈ [0, 1)
		// The marginal distributions of tasks’ delays.
		Marginal string
		// The strength of correlations between tasks.
		CorrLength float64 // > 0
		// The portion of the variance to be preserved when reducing the number
		// of stochastic dimensions.
		VarThreshold float64 // ∈ (0, 1]
	}

	// The quantities of interest are “end-to-end-delay,” “total-energy,” and
	// “temperature-profile.”
	Target string

	// The configuration of temperature analysis. Specific to the
	// temperature-profile target.
	Temperature struct {
		// The indices of the cores that should be considered; if empty, the
		// parameter is set to all cores.
		Cores []uint
		// The time step.
		TimeStep float64
		// The indices of the time moments that should be considered; if empty,
		// the parameter is set to all time moments.
		Steps []uint
		numeric.Config
	}

	// The configuration of interpolation.
	Interpolation struct {
		Rule string
		adhier.Config
	}

	Assessment struct {
		// The seed for random number generation.
		Seed int
		// The number of samples to take.
		Samples uint
	}

	// A flag for displaying progress information.
	Verbose bool
}

func NewConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var config Config
	if err = decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
