package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"runtime"
	"strconv"

	"github.com/simulated-reality/laboratory/cmd/internal"
	"github.com/simulated-reality/laboratory/internal/config"
)

var (
	outputFile  = flag.String("o", "", "an output file (required)")
	sampleSeed  = flag.String("s", "", "a seed for generating samples")
	sampleCount = flag.String("n", "", "the number of samples")
)

type Config *config.Assessment

func main() {
	internal.Run(command)
}

func command(globalConfig *config.Config) error {
	globalConfig.Probability.VarThreshold = math.Inf(1)

	config := &globalConfig.Assessment
	if len(*sampleSeed) > 0 {
		if number, err := strconv.ParseInt(*sampleSeed, 0, 64); err != nil {
			return err
		} else {
			config.Seed = number
		}
	}
	if len(*sampleCount) > 0 {
		if number, err := strconv.ParseUint(*sampleCount, 0, 64); err != nil {
			return err
		} else {
			config.Samples = uint(number)
		}
	}

	if config.Samples == 0 {
		return errors.New("the number of samples should be positive")
	}

	output, err := internal.Create(*outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	problem, err := internal.NewProblem(globalConfig)
	if err != nil {
		return err
	}

	target, err := internal.NewTarget(problem)
	if err != nil {
		return err
	}

	ni, no := target.Dimensions()
	ns := config.Samples

	points := internal.Generate(ni, ns, config.Seed)

	if globalConfig.Verbose {
		fmt.Printf("Evaluating the original model at %d points...\n", ns)
	}

	values := internal.Invoke(target, points, uint(runtime.GOMAXPROCS(0)))

	if globalConfig.Verbose {
		fmt.Println("Done.")
	}

	if err := output.Put("points", points, ni, ns); err != nil {
		return err
	}
	if err := output.Put("values", values, no, ns); err != nil {
		return err
	}

	return nil
}
