package internal

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/ready-steady/format/mat"
)

func Run(command func(*Problem, *mat.File, *mat.File) error) {
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

	if len(*configFile) == 0 {
		printError(errors.New("a configuration file is required"))
		return
	}

	problem, err := NewProblem(*configFile)
	if err != nil {
		printError(err)
		return
	}

	var ifile, ofile *mat.File

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

	if err = command(problem, ifile, ofile); err != nil {
		printError(err)
		return
	}
}

func printError(err error) {
	fmt.Printf("Error: %s.\n\n", err)

	fmt.Printf("Usage: %s [options]", os.Args[0])
	fmt.Printf(`

Options:
    -c <FILE.json>  - a configuration file (required)
    -i <FILE.mat>   - an input file
    -o <FILE.mat>   - an output file
`)
}
