package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/ready-steady/format/mat"
)

var config = flag.String("c", "", "")
var input = flag.String("i", "", "")
var output = flag.String("o", "", "")

func main() {
	if len(os.Args) == 1 {
		printUsage()
		return
	}

	var command func(*problem, *mat.File, *mat.File) error
	var problem *problem
	var ifile, ofile *mat.File

	if command = findCommand(os.Args[1]); command == nil {
		printError(errors.New("the command is unknown"))
		return
	}

	// Remove the command name.
	os.Args[1] = os.Args[0]
	os.Args = os.Args[1:]

	flag.Parse()

	if len(*config) == 0 {
		printError(errors.New("a problem specification is required"))
		return
	}

	config, err := loadConfig(*config)
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

	if len(*input) > 0 {
		if ifile, err = mat.Open(*input, "r"); err != nil {
			printError(err)
			return
		}
		defer ifile.Close()
	}

	if len(*output) > 0 {
		if ofile, err = mat.Open(*output, "w7.3"); err != nil {
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

func printUsage() {
	fmt.Printf("Usage: %s <command> [options]", os.Args[0])
	fmt.Printf(`

Commands:
    show        - to display the configuration of a problem
    solve       - to construct a surrogate model
    check       - to assess a surrogate model

Options:
    -c <FILE>   - a problem specification in JSON (required)
    -i <FILE>   - an input MAT file
    -o <FILE>   - an output MAT file
`)
}

func printError(err error) {
	fmt.Printf("Error: %s.\n", err)
}
