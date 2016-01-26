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
	"github.com/turing-complete/laboratory/src/internal/support"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/target"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"

	isolver "github.com/turing-complete/laboratory/src/internal/solver"
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

	auncertainty, err := uncertainty.NewAleatory(system, &config.Uncertainty)
	if err != nil {
		return err
	}

	euncertainty, err := uncertainty.NewEpistemic(system, &config.Uncertainty)
	if err != nil {
		return err
	}

	atarget, err := target.New(system, auncertainty, &config.Target)
	if err != nil {
		return err
	}

	etarget, err := target.New(system, euncertainty, &config.Target)
	if err != nil {
		return err
	}

	ni, no := etarget.Dimensions()

	solver, err := isolver.New(ni, no, &config.Solver)
	if err != nil {
		return err
	}

	solution := new(isolver.Solution)
	if err = approximate.Get("solution", solution); err != nil {
		return err
	}

	ns := config.Assessment.Samples

	points := generate(etarget, atarget, ns, config.Assessment.Seed)

	log.Printf("Evaluating the surrogate model at %d points...\n", ns)
	log.Printf("%10s %15s\n", "Iteration", "Nodes")

	nk := uint(len(solution.Active))

	steps := make([]uint, nk)
	values := make([]float64, 0, ns*no)

	k, Δ := uint(0), float64(nk-1)/(math.Min(maxSteps, float64(nk))-1)

	for i, na := uint(0), uint(0); i < nk; i++ {
		na += solution.Active[i]
		steps[k] += solution.Active[i]

		if i != uint(float64(k)*Δ+0.5) {
			continue
		}
		k++

		log.Printf("%10d %15d\n", i, na)

		s := *solution
		s.Nodes = na
		s.Indices = s.Indices[:na*ni]
		s.Surpluses = s.Surpluses[:na*no]

		values = append(values, solver.Evaluate(&s, points)...)
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

	return nil
}

func generate(into, from target.Target, ns uint, seed int64) []float64 {
	nif, _ := from.Dimensions()
	nii, _ := into.Dimensions()

	zf := support.Generate(nif, ns, seed)
	zi := make([]float64, nii*ns)

	for i := uint(0); i < ns; i++ {
		copy(zi[i*nii:(i+1)*nii], into.Forward(from.Inverse(zf[i*nif:(i+1)*nif])))
	}

	return zi
}
