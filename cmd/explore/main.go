package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"runtime"
	"strings"

	"github.com/ready-steady/linear"

	"../internal"
)

var (
	outputFile     = flag.String("o", "", "an output file (required)")
	parameterIndex = flag.String("s", "[]", "the parameters to sweep")
	numberOfNodes  = flag.Uint("n", 10, "the number of nodes per parameter")
	defaultNode    = flag.Float64("d", 0.5, "the default value of parameters")
)

func main() {
	internal.Run(command)
}

func command(config *internal.Config) error {
	output, err := internal.Create(*outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	problem, err := internal.NewProblem(config)
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
	nn := *numberOfNodes

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
			return nil, errors.New(fmt.Sprintf("the indices should be less that %v", ni))
		}
	}

	return index, nil
}
