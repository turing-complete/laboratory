package main

import (
	"flag"
	"log"

	"github.com/turing-complete/laboratory/src/internal/command"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/database"
	"github.com/turing-complete/laboratory/src/internal/quantity"
	"github.com/turing-complete/laboratory/src/internal/solution"
	"github.com/turing-complete/laboratory/src/internal/system"
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

	auncertainty, err := uncertainty.NewAleatory(system, &config.Uncertainty)
	if err != nil {
		return err
	}
	aquantity, err := quantity.New(system, auncertainty, &config.Quantity)
	if err != nil {
		return err
	}

	euncertainty, err := uncertainty.NewEpistemic(system, &config.Uncertainty)
	if err != nil {
		return err
	}
	equantity, err := quantity.New(system, euncertainty, &config.Quantity)
	if err != nil {
		return err
	}

	var target, reference quantity.Quantity
	if config.Solution.Aleatory {
		target, reference = aquantity, equantity
	} else {
		target, reference = equantity, aquantity
	}

	ni, no := target.Dimensions()

	solution, err := solution.New(ni, no, &config.Solution)
	if err != nil {
		return err
	}

	log.Println("System", system)
	log.Println("Quantity", target)

	if config.Solution.Aleatory {
		log.Println("Constructing an aleatory surrogate...")
	} else {
		log.Println("Constructing an epistemic surrogate...")
	}

	surrogate := solution.Compute(target, reference)

	log.Println("Surrogate", surrogate)

	if err := output.Put("surrogate", *surrogate); err != nil {
		return err
	}

	return nil
}
