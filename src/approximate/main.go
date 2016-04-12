package main

import (
	"flag"
	"log"

	"github.com/turing-complete/laboratory/src/internal/command"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/database"
	"github.com/turing-complete/laboratory/src/internal/solver"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/target"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
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

	system, err := system.New(&config.System)
	if err != nil {
		return err
	}

	uncertainty, err := uncertainty.NewEpistemic(system, &config.Uncertainty)
	if err != nil {
		return err
	}

	target, err := target.New(system, uncertainty, &config.Target)
	if err != nil {
		return err
	}

	ni, no := target.Dimensions()

	solver, err := solver.New(ni, no, &config.Solver)
	if err != nil {
		return err
	}

	log.Println("System", system)
	log.Println("Target", target)
	log.Println("Constructing a surrogate...")

	solution := solver.Compute(target)
	log.Println("Solution", solution)

	if err := output.Put("solution", *solution); err != nil {
		return err
	}

	return nil
}
