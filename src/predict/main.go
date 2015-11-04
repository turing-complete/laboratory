package main

import (
	"errors"
	"flag"
	"log"
	"math"
	"strconv"

	"github.com/turing-complete/laboratory/src/internal/command"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/database"
	"github.com/turing-complete/laboratory/src/internal/solver"
	"github.com/turing-complete/laboratory/src/internal/support"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/target"
)

var (
	approximateFile = flag.String("approximate", "", "an output of `approximate` (required)")
	outputFile      = flag.String("o", "", "an output file (required)")
	sampleSeed      = flag.String("s", "", "a seed for generating samples")
	sampleCount     = flag.String("n", "", "the number of samples")
)

type Config *config.Assessment

func main() {
	command.Run(function)
}

func function(config *config.Config) error {
	const (
		maxSteps = 10
	)

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

	approximate, err := database.Open(*approximateFile)
	if err != nil {
		return err
	}
	defer approximate.Close()

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

	aSolver, err := solver.New(target, &config.Solver)
	if err != nil {
		return err
	}

	solution := new(solver.Solution)
	if err = approximate.Get("solution", solution); err != nil {
		return err
	}

	ni, no := target.Dimensions()
	ns := config.Assessment.Samples

	points := support.Generate(ni, ns, config.Assessment.Seed)

	log.Printf("Evaluating the surrogate model at %d points...\n", ns)
	log.Printf("%10s %15s %15s\n", "Iteration", "Accepted Nodes", "Rejected Nodes")

	nk := uint(len(solution.Accept))

	steps := make([]uint, nk)
	values := make([]float64, 0, ns*no)
	moments := make([]float64, 0, no)

	k, Δ := uint(0), float64(nk-1)/(math.Min(maxSteps, float64(nk))-1)

	for i, na, nr := uint(0), uint(0), uint(0); i < nk; i++ {
		na += solution.Accept[i]
		nr += solution.Reject[i]

		steps[k] += solution.Accept[i] + solution.Reject[i]

		if i != uint(float64(k)*Δ+0.5) {
			continue
		}
		k++

		log.Printf("%10d %15d %15d\n", i, na, nr)

		s := *solution
		s.Nodes = na
		s.Indices = s.Indices[:na*ni]
		s.Surpluses = s.Surpluses[:na*no]

		values = append(values, aSolver.Evaluate(&s, points)...)
		moments = append(moments, aSolver.Integrate(&s)...)
	}

	nk, steps = k, steps[:k]

	log.Println("Done.")

	if err := output.Put("solution", *solution); err != nil {
		return err
	}
	if err := output.Put("points", points, ni, ns); err != nil {
		return err
	}
	if err := output.Put("steps", steps); err != nil {
		return err
	}
	if err := output.Put("values", values, no, ns, nk); err != nil {
		return err
	}
	if err := output.Put("moments", moments, no, nk); err != nil {
		return err
	}

	return nil
}
