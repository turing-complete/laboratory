package internal

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/ready-steady/format/mat"
	"github.com/ready-steady/numeric/basis/linhat"
	"github.com/ready-steady/numeric/grid/newcot"
	"github.com/ready-steady/numeric/interpolation/adhier"
)

func Run(command func(*Config, *Problem, *mat.File, *mat.File) error) {
	configFile := flag.String("c", "", "")
	inputFile := flag.String("i", "", "")
	outputFile := flag.String("o", "", "")

	profile := flag.String("profile", "", "")

	flag.Parse()

	if len(*profile) > 0 {
		pfile, err := os.Create(*profile)
		if err != nil {
			printError(errors.New("cannot enable profiling"))
			return
		}
		pprof.StartCPUProfile(pfile)
		defer pprof.StopCPUProfile()
	}

	var problem *Problem
	var ifile, ofile *mat.File

	if len(*configFile) == 0 {
		printError(errors.New("a problem specification is required"))
		return
	}

	config, err := loadConfig(*configFile)
	if err != nil {
		printError(err)
		return
	}

	if err = config.validate(); err != nil {
		printError(err)
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	if problem, err = newProblem(&config); err != nil {
		printError(err)
		return
	}

	if len(*inputFile) > 0 {
		if ifile, err = mat.Open(*inputFile, "r"); err != nil {
			printError(err)
			return
		}
		defer ifile.Close()
	}

	if len(*outputFile) > 0 {
		if ofile, err = mat.Open(*outputFile, "w7.3"); err != nil {
			printError(err)
			return
		}
		defer ofile.Close()
	}

	if err = command(&config, problem, ifile, ofile); err != nil {
		printError(err)
		return
	}
}

func Setup(p *Problem) (Target, *adhier.Interpolator, error) {
	c := p.config

	var target Target
	var err error

	switch p.config.Target {
	case "end-to-end-delay":
		target = newDelayTarget(p)
	case "total-energy":
		target = newEnergyTarget(p)
	case "temperature-profile":
		target, err = newTemperatureTarget(p)
	default:
		err = errors.New("the target is unknown")
	}
	if err != nil {
		return nil, nil, err
	}

	ic, oc := target.InputsOutputs()

	var grid adhier.Grid
	var basis adhier.Basis

	switch strings.ToLower(c.Interpolation.Rule) {
	case "open":
		grid = newcot.NewOpen(uint16(ic))
		basis = linhat.NewOpen(uint16(ic), uint16(oc))
	case "closed":
		grid = newcot.NewClosed(uint16(ic))
		basis = linhat.NewClosed(uint16(ic), uint16(oc))
	default:
		return nil, nil, errors.New("the interpolation rule is unknown")
	}

	interpolator, err := adhier.New(grid, basis, adhier.Config(c.Interpolation.Config))
	if err != nil {
		return nil, nil, err
	}

	return target, interpolator, nil
}

func printError(err error) {
	fmt.Printf("Error: %s.\n\n", err)

	fmt.Printf("Usage: %s [options]", os.Args[0])
	fmt.Printf(`

Options:
    -c <FILE>   - a problem specification in JSON (required)
    -i <FILE>   - an input MAT file
    -o <FILE>   - an output MAT file
`)
}
