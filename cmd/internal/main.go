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
	profileFile := flag.String("p", "", "")

	flag.Parse()

	if len(*profileFile) > 0 {
		pfile, err := os.Create(*profileFile)
		if err != nil {
			printError(errors.New("cannot enable profiling"))
			return
		}
		pprof.StartCPUProfile(pfile)
		defer pprof.StopCPUProfile()
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	var ifile, ofile *mat.File

	problem, err := SetupProblem(*configFile)
	if err != nil {
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

	if err = command(problem.config, problem, ifile, ofile); err != nil {
		printError(err)
		return
	}
}

func SetupProblem(configFile string) (*Problem, error) {
	if len(configFile) == 0 {
		return nil, errors.New("a problem specification is required")
	}

	config, err := loadConfig(configFile)
	if err != nil {
		return nil, err
	}
	if err = config.validate(); err != nil {
		return nil, err
	}

	return newProblem(config)
}

func SetupTarget(problem *Problem) (Target, error) {
	switch problem.config.Target {
	case "end-to-end-delay":
		return newDelayTarget(problem), nil
	case "total-energy":
		return newEnergyTarget(problem), nil
	case "temperature-profile":
		return newTemperatureTarget(problem)
	default:
		return nil, errors.New("the target is unknown")
	}
}

func SetupInterpolator(problem *Problem, target Target) (*adhier.Interpolator, error) {
	config := &problem.config.Interpolation
	ic, oc := target.InputsOutputs()

	var grid adhier.Grid
	var basis adhier.Basis

	switch strings.ToLower(config.Rule) {
	case "open":
		grid = newcot.NewOpen(uint16(ic))
		basis = linhat.NewOpen(uint16(ic), uint16(oc))
	case "closed":
		grid = newcot.NewClosed(uint16(ic))
		basis = linhat.NewClosed(uint16(ic), uint16(oc))
	default:
		return nil, errors.New("the interpolation rule is unknown")
	}

	return adhier.New(grid, basis, (*adhier.Config)(&config.Config))
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
