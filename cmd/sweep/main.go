package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math"
	"runtime"
	"strings"

	"github.com/ready-steady/linear"
	"github.com/simulated-reality/laboratory/cmd/internal"
	"github.com/simulated-reality/laboratory/internal/config"
	"github.com/simulated-reality/laboratory/internal/file"
	"github.com/simulated-reality/laboratory/internal/problem"
)

var (
	outputFile     = flag.String("o", "", "an output file (required)")
	varThreshold   = flag.Float64("t", math.Inf(1), "the variance-reduction threshold")
	parameterIndex = flag.String("s", "[]", "the parameters to sweep")
	defaultNode    = flag.Float64("d", 0.5, "the default value of parameters")
	nodeCount      = flag.Uint("n", 10, "the number of nodes per parameter")
)

func main() {
	internal.Run(command)
}

func command(config *config.Config) error {
	output, err := file.Create(*outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	config.Probability.VarThreshold = *varThreshold

	problem, err := problem.New(config)
	if err != nil {
		return err
	}

	target, err := internal.NewTarget(problem)
	if err != nil {
		return err
	}

	points, err := generate(target, config.Interpolation.Rule)
	if err != nil {
		return err
	}

	ni, no := target.Dimensions()
	np := uint(len(points)) / ni

	if config.Verbose {
		fmt.Println(problem)
		fmt.Println(target)
		fmt.Printf("Evaluating the model with reduction %.2f at %v points...\n",
			config.Probability.VarThreshold, np)
	}

	values := internal.Invoke(target, points, uint(runtime.GOMAXPROCS(0)))

	if config.Verbose {
		fmt.Println("Done.")
	}

	if err := output.Put("values", values, no, np); err != nil {
		return err
	}
	if err := output.Put("points", points, ni, np); err != nil {
		return err
	}

	return nil
}

func generate(target internal.Target, rule string) ([]float64, error) {
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

func detect(target internal.Target) ([]uint, error) {
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
