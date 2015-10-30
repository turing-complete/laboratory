package main

import (
	"flag"
	"fmt"

	"github.com/turing-complete/laboratory/src/internal/command"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/database"
	"github.com/turing-complete/laboratory/src/internal/problem"
	"github.com/turing-complete/laboratory/src/internal/solver"
	"github.com/turing-complete/laboratory/src/internal/target"
)

var (
	outputFile = flag.String("o", "", "an output file (required)")
)

func main() {
	command.Run(function)
}

func function(config *config.Config) error {
	output, err := database.Create(*outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	problem, err := problem.New(config)
	if err != nil {
		return err
	}

	target, err := target.New(problem, &config.Target)
	if err != nil {
		return err
	}

	solver, err := solver.New(problem, target, &config.Interpolation)
	if err != nil {
		return err
	}

	if config.Verbose {
		fmt.Println(problem)
		fmt.Println(target)
		fmt.Println("Constructing a surrogate...")
	}

	solution := solver.Compute(target)
	if config.Verbose {
		fmt.Println(solution)
	}
	if err := output.Put("solution", *solution); err != nil {
		return err
	}

	return nil
}
