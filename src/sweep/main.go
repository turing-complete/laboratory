package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/ready-steady/linear"
	"github.com/turing-complete/laboratory/src/internal/command"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/database"
	"github.com/turing-complete/laboratory/src/internal/solution"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/target"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

var (
	approximateFile = flag.String("approximate", "", "an output of `approximate`")
	outputFile      = flag.String("o", "", "an output file (required)")
	parameterIndex  = flag.String("s", "[]", "the parameters to sweep")
	defaultNode     = flag.Float64("d", 0.5, "the default value of parameters")
	nodeCount       = flag.Uint("n", 10, "the number of nodes per parameter")
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

	atarget, err := target.New(system, uncertainty, &config.Target)
	if err != nil {
		return err
	}

	points, err := generate(atarget, config.Solution.Rule)
	if err != nil {
		return err
	}

	ni, no := atarget.Dimensions()
	np := uint(len(points)) / ni

	log.Println(system)
	log.Println(atarget)

	var values []float64
	if len(*approximateFile) > 0 {
		approximate, err := database.Open(*approximateFile)
		if err != nil {
			return err
		}
		defer approximate.Close()

		asolution, err := solution.New(ni, no, &config.Solution)
		if err != nil {
			return err
		}

		surrogate := new(solution.Surrogate)
		if err = approximate.Get("surrogate", surrogate); err != nil {
			return err
		}

		log.Printf("Evaluating the approximation at %d points...\n", np)
		values = asolution.Evaluate(surrogate, points)
	} else {
		log.Printf("Evaluating the original model at %d points...\n", np)
		values = target.Invoke(atarget, points)
	}

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

	return linear.TensorFloat64(parameters...), nil
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
