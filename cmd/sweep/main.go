package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"runtime"
	"strings"

	"github.com/ready-steady/hdf5"
	"github.com/ready-steady/linear"

	"../internal"
)

var (
	parameters = flag.String("s", "[]", "the parameters to sweep")
	pointCount = flag.Uint("n", 10, "the number of points per parameter")
)

func main() {
	internal.Run(command)
}

func command(config internal.Config, _ *hdf5.File, output *hdf5.File) error {
	problem, err := internal.NewProblem(config)
	if err != nil {
		return err
	}

	target, err := internal.NewTarget(problem)
	if err != nil {
		return err
	}

	points, err := generate(target)
	if err != nil {
		return err
	}

	ni, no := target.Dimensions()
	np := uint(len(points)) / ni

	if config.Verbose {
		fmt.Printf("Evaluating the reduced model at %v points...\n", np)
		fmt.Println(problem)
		fmt.Println(target)
	}

	values := internal.Invoke(target, points, uint(runtime.GOMAXPROCS(0)))

	if config.Verbose {
		fmt.Println("Done.")
	}

	if output == nil {
		return nil
	}

	if err := output.Put("values", values, no, np); err != nil {
		return err
	}
	if err := output.Put("points", points, ni, np); err != nil {
		return err
	}

	return nil
}

func generate(target internal.Target) ([]float64, error) {
	ni, _ := target.Dimensions()
	np := *pointCount

	index, err := detect(target)
	if err != nil {
		return nil, err
	}

	parameters := make([][]float64, ni)

	steady := []float64{0.5}
	for i := uint(0); i < ni; i++ {
		parameters[i] = steady
	}

	sweep := make([]float64, np)
	for i := uint(0); i < np; i++ {
		sweep[i] = float64(i) * 1.0 / float64(np-1)
	}
	for _, i := range index {
		parameters[i] = sweep
	}

	return linear.Tensor(parameters...), nil
}

func detect(target internal.Target) ([]uint, error) {
	ni, _ := target.Dimensions()

	index := []uint{}

	decoder := json.NewDecoder(strings.NewReader(*parameters))
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
