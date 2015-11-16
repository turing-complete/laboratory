package main

import (
	"errors"
	"flag"
	"log"
	"math"
	"runtime"
	"strconv"

	"github.com/turing-complete/laboratory/src/internal/command"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/database"
	"github.com/turing-complete/laboratory/src/internal/support"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/target"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"

	itarget "github.com/turing-complete/laboratory/src/internal/target"
)

var (
	outputFile  = flag.String("o", "", "an output file (required)")
	sampleSeed  = flag.String("s", "", "a seed for generating samples")
	sampleCount = flag.String("n", "", "the number of samples")
)

type Config *config.Assessment

func main() {
	command.Run(function)
}

func function(config *config.Config) error {
	config.Uncertainty.Reduction = math.Inf(1)

	if len(*sampleSeed) > 0 {
		if number, err := strconv.ParseInt(*sampleSeed, 0, 64); err != nil {
			return err
		} else {
			config.Assessment.Seed = number
		}
	}
	if len(*sampleCount) > 0 {
		if number, err := strconv.ParseUint(*sampleCount, 0, 64); err != nil {
			return err
		} else {
			config.Assessment.Samples = uint(number)
		}
	}

	if config.Assessment.Samples == 0 {
		return errors.New("the number of samples should be positive")
	}

	output, err := database.Create(*outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	system, err := system.New(&config.System)
	if err != nil {
		return err
	}

	uncertainty, err := uncertainty.NewMarginal(system, &config.Uncertainty)
	if err != nil {
		return err
	}

	target, err := target.New(system, uncertainty, &config.Target)
	if err != nil {
		return err
	}

	ni, no := target.Dimensions()
	ns := config.Assessment.Samples

	points := support.Generate(ni, ns, config.Assessment.Seed)

	log.Printf("Evaluating the original model at %d points...\n", ns)
	values := itarget.Invoke(target, points, uint(runtime.GOMAXPROCS(0)))
	log.Println("Done.")

	if err := output.Put("points", points, ni, ns); err != nil {
		return err
	}
	if err := output.Put("values", values, no, ns); err != nil {
		return err
	}

	return nil
}
