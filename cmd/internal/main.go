package internal

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/ready-steady/format/mat"
)

func Run(command func(*Config, *Problem, *mat.File, *mat.File) error) {
	configFile := flag.String("c", "", "")
	inputFile := flag.String("i", "", "")
	outputFile := flag.String("o", "", "")

	flag.Parse()

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

	if problem, err = newProblem(config); err != nil {
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

	if err = command(&problem.config, problem, ifile, ofile); err != nil {
		printError(err)
		return
	}
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
