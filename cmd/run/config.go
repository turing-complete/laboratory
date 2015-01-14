package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/ready-steady/tempan/expint"

	"../../pkg/solver"
)

type Config struct {
	// The TGFF file of the system to analyze.
	TGFF string

	// The IDs of the cores to analyze; if empty, set to all cores.
	CoreIndex []uint16
	// The IDs of the tasks to analyze; if empty, set to all tasks.
	TaskIndex []uint16

	ProbModel struct {
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

	// The quantity of interest. Available targes are “end-to-end-delay” and
	// “temperature-profile.”
	Target string

	// The configuration of the algorithm for temperature analysis. Specific to
	// the temperature-profile target.
	TempAnalysis expint.Config

	// The configuration of the solver.
	Solver solver.Config

	// The seed for random number generation.
	Seed int64
	// The number of samples to take.
	Samples uint32

	// A flag for displaying progress information.
	Verbose bool
}

func loadConfig(path string) (Config, error) {
	c := Config{}

	file, err := os.Open(path)
	if err != nil {
		return c, err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	if err = dec.Decode(&c); err != nil {
		return c, err
	}

	return c, nil
}

func (c *Config) validate() error {
	if c.ProbModel.MaxDelay < 0 || 1 <= c.ProbModel.MaxDelay {
		return errors.New("the delay rate is invalid")
	}
	if c.ProbModel.CorrLength <= 0 {
		return errors.New("the correlation length is invalid")
	}
	if c.ProbModel.VarThreshold <= 0 || 1 < c.ProbModel.VarThreshold {
		return errors.New("the variance-reduction threshold is invalid")
	}

	return nil
}
