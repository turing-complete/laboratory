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

	target, err := target.New(system, &config.Target)
	if err != nil {
		return err
	}

	solver, err := solver.New(target, &config.Solver)
	if err != nil {
		return err
	}

	log.Println(system)
	log.Println(target)
	log.Println("Constructing a surrogate...")

	solution := solver.Compute(target)
	log.Println(solution)

	if err := output.Put("solution", *solution); err != nil {
		return err
	}

	return nil
}
