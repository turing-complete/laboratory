package main

import (
	"errors"
	"flag"
	"log"
	"strconv"

	"github.com/ready-steady/choose"
	"github.com/turing-complete/laboratory/src/internal/command"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/database"
	"github.com/turing-complete/laboratory/src/internal/quantity"
	"github.com/turing-complete/laboratory/src/internal/solution"
	"github.com/turing-complete/laboratory/src/internal/support"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

const (
	maxSteps = 10
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
	aquantity, err := quantity.New(system, auncertainty, &config.Quantity)
	if err != nil {
		return err
	}

	euncertainty, err := uncertainty.NewEpistemic(system, &config.Uncertainty)
	if err != nil {
		return err
	}
	equantity, err := quantity.New(system, euncertainty, &config.Quantity)
	if err != nil {
		return err
	}

	var target, proxy quantity.Quantity
	if config.Solution.Aleatory {
		target, proxy = aquantity, aquantity // noop
	} else {
		target, proxy = equantity, aquantity
	}

	ni, no := target.Dimensions()
	ns := config.Assessment.Samples

	surrogate := new(solution.Surrogate)
	if err = approximate.Get("surrogate", surrogate); err != nil {
		return err
	}

	solution, err := solution.New(ni, no, &config.Solution)
	if err != nil {
		return err
	}

	points := generate(target, proxy, ns, config.Assessment.Seed)

	log.Printf("Evaluating the surrogate model at %d points...\n", ns)
	log.Printf("%5s %15s\n", "Step", "Nodes")

	nk := uint(len(surrogate.Active))

	cumsum := append([]uint(nil), surrogate.Active...)
	for i := uint(1); i < nk; i++ {
		cumsum[i] += cumsum[i-1]
	}
	indices := choose.UniformUint(cumsum, maxSteps)

	nk = uint(len(indices))

	active := make([]uint, nk)
	for i := uint(0); i < nk; i++ {
		active[i] = cumsum[indices[i]]
	}

	values := make([]float64, 0, ns*no)
	for i := uint(0); i < nk; i++ {
		log.Printf("%5d %15d\n", i, active[i])

		s := *surrogate
		s.Nodes = active[i]
		s.Indices = s.Indices[:active[i]*ni]
		s.Surpluses = s.Surpluses[:active[i]*no]

		if !solution.Validate(&s) {
			panic("something went wrong")
		}

		values = append(values, solution.Evaluate(&s, points)...)
	}

	if err := output.Put("surrogate", *surrogate); err != nil {
		return err
	}
	if err := output.Put("points", points, ni, ns); err != nil {
		return err
	}
	if err := output.Put("values", values, no, ns, nk); err != nil {
		return err
	}
	if err := output.Put("active", active); err != nil {
		return err
	}

	return nil
}

func generate(into, from quantity.Quantity, ns uint, seed int64) []float64 {
	ni, _ := into.Dimensions()
	nf, _ := from.Dimensions()
	zi := make([]float64, ni*ns)
	zf := support.Generate(nf, ns, seed)
	for i := uint(0); i < ns; i++ {
		copy(zi[i*ni:(i+1)*ni], into.Forward(from.Backward(zf[i*nf:(i+1)*nf])))
	}
	return zi
}
