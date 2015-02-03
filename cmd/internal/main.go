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

func Run(command func(string, *mat.File, *mat.File) error) {
	config := flag.String("c", "", "")
	input := flag.String("i", "", "")
	output := flag.String("o", "", "")
	profile := flag.String("p", "", "")

	flag.Parse()

	var err error

	if len(*profile) > 0 {
		pfile, err := os.Create(*profile)
		if err != nil {
			printError(errors.New("cannot enable profiling"))
			return
		}
		pprof.StartCPUProfile(pfile)
		defer pprof.StopCPUProfile()
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	var ifile, ofile *mat.File

	if len(*config) == 0 {
		printError(errors.New("a configuration file is required"))
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

	if err := command(*config, ifile, ofile); err != nil {
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
