package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"strconv"

	"github.com/simulated-reality/laboratory/cmd/internal"
	"github.com/simulated-reality/laboratory/internal/config"
	"github.com/simulated-reality/laboratory/internal/database"
	"github.com/simulated-reality/laboratory/internal/problem"
	"github.com/simulated-reality/laboratory/internal/solver"
	"github.com/simulated-reality/laboratory/internal/support"
	"github.com/simulated-reality/laboratory/internal/target"
)

var (
	approximateFile = flag.String("approximate", "", "an output of `approximate` (required)")
	outputFile      = flag.String("o", "", "an output file (required)")
	sampleSeed      = flag.String("s", "", "a seed for generating samples")
	sampleCount     = flag.String("n", "", "the number of samples")
)

type Config *config.Assessment

func main() {
	internal.Run(command)
}

func command(globalConfig *config.Config) error {
	const (
		maxSteps = 10
	)

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

	problem, err := problem.New(globalConfig)
	if err != nil {
		return err
	}

	target, err := target.New(problem)
	if err != nil {
		return err
	}

	aSolver, err := solver.New(problem, target)
	if err != nil {
		return err
	}

	solution := new(solver.Solution)
	if err = approximate.Get("solution", solution); err != nil {
		return err
	}

	ni, no := target.Dimensions()
	ns := config.Samples

	points := support.Generate(ni, ns, config.Seed)

	if globalConfig.Verbose {
		fmt.Printf("Evaluating the surrogate model at %d points...\n", ns)
		fmt.Printf("%10s %15s %15s\n", "Iteration", "Accepted Nodes", "Rejected Nodes")
	}

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

		if globalConfig.Verbose {
			fmt.Printf("%10d %15d %15d\n", i, na, nr)
		}

		s := *solution
		s.Nodes = na
		s.Indices = s.Indices[:na*ni]
		s.Surpluses = s.Surpluses[:na*no]

		values = append(values, aSolver.Evaluate(&s, points)...)
		moments = append(moments, aSolver.Integrate(&s)...)
	}

	nk, steps = k, steps[:k]

	if globalConfig.Verbose {
		fmt.Println("Done.")
	}

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
