package main

import (
	"flag"
	"fmt"

	"github.com/simulated-reality/laboratory/cmd/internal"
	"github.com/simulated-reality/laboratory/internal/config"
	"github.com/simulated-reality/laboratory/internal/database"
	"github.com/simulated-reality/laboratory/internal/problem"
	"github.com/simulated-reality/laboratory/internal/solver"
	"github.com/simulated-reality/laboratory/internal/target"
)

var (
	outputFile = flag.String("o", "", "an output file (required)")
)

func main() {
	internal.Run(command)
}

func command(config *config.Config) error {
	output, err := database.Create(*outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	problem, err := problem.New(config)
	if err != nil {
		return err
	}

	target, err := target.New(problem)
	if err != nil {
		return err
	}

	solver, err := solver.New(problem, target)
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
