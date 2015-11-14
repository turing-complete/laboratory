package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"runtime"
	"strings"

	"github.com/ready-steady/linear"
	"github.com/turing-complete/laboratory/src/internal/command"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/database"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/target"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

var (
	outputFile     = flag.String("o", "", "an output file (required)")
	reduction      = flag.Float64("r", math.Inf(1), "the variance reduction")
	parameterIndex = flag.String("s", "[]", "the parameters to sweep")
	defaultNode    = flag.Float64("d", 0.5, "the default value of parameters")
	nodeCount      = flag.Uint("n", 10, "the number of nodes per parameter")
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

	config.Uncertainty.Reduction = *reduction

	system, err := system.New(&config.System)
	if err != nil {
		return err
	}

	uncertainty, err := uncertainty.New(system, &config.Uncertainty)
	if err != nil {
		return err
	}

	aTarget, err := target.New(system, uncertainty, &config.Target)
	if err != nil {
		return err
	}

	points, err := generate(aTarget, config.Solver.Rule)
	if err != nil {
		return err
	}

	ni, no := aTarget.Dimensions()
	np := uint(len(points)) / ni

	log.Println(system)
	log.Println(aTarget)
	log.Printf("Evaluating the model with reduction %.2f at %v points...\n",
		config.Uncertainty.Reduction, np)

	values := target.Invoke(aTarget, points, uint(runtime.GOMAXPROCS(0)))

	log.Println("Done.")

	if err := output.Put("values", values, no, np); err != nil {
		return err
	}
	if err := output.Put("points", points, ni, np); err != nil {
		return err
	}

	return nil
}

func generate(target target.Target, rule string) ([]float64, error) {
	ni, _ := target.Dimensions()
	nn := *nodeCount

	index, err := detect(target)
	if err != nil {
		return nil, err
	}

	steady := []float64{*defaultNode}

	sweep := make([]float64, nn)
	switch rule {
	case "closed":
		for i := uint(0); i < nn; i++ {
			sweep[i] = float64(i) / float64(nn-1)
		}
	case "open":
		for i := uint(0); i < nn; i++ {
			sweep[i] = float64(i+1) / float64(nn+1)
		}
	default:
		return nil, errors.New("the sweep rule is unknown")
	}

	parameters := make([][]float64, ni)
	for i := uint(0); i < ni; i++ {
		parameters[i] = steady
	}
	for _, i := range index {
		parameters[i] = sweep
	}

	return linear.Tensor(parameters...), nil
}

func detect(target target.Target) ([]uint, error) {
	ni, _ := target.Dimensions()

	index := []uint{}

	decoder := json.NewDecoder(strings.NewReader(*parameterIndex))
	if err := decoder.Decode(&index); err != nil {
		return nil, err
	}

	if len(index) == 0 {
		index = make([]uint, ni)
		for i := uint(0); i < ni; i++ {
			index[i] = i
		}
	}

	for _, i := range index {
		if i >= ni {
			return nil, errors.New(fmt.Sprintf("the indices should be less than %v", ni))
		}
	}

	return index, nil
}
