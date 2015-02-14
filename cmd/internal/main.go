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

func Run(command func(Config, *mat.File, *mat.File) error) {
	configFile := flag.String("c", "", "")
	inputFile := flag.String("i", "", "")
	outputFile := flag.String("o", "", "")
	profileFile := flag.String("p", "", "")

	flag.Parse()

	var err error

	if len(*profileFile) > 0 {
		profile, err := os.Create(*profileFile)
		if err != nil {
			fail(errors.New("cannot enable profiling"))
		}
		pprof.StartCPUProfile(profile)
		defer pprof.StopCPUProfile()
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	var input, output *mat.File

	if len(*configFile) == 0 {
		fail(errors.New("a configuration file is required"))
	}
	config, err := NewConfig(*configFile)
	if err != nil {
		fail(err)
	}

	if len(*inputFile) > 0 {
		if input, err = mat.Open(*inputFile, "r"); err != nil {
			fail(err)
		}
		defer input.Close()
	}

	if len(*outputFile) > 0 {
		if output, err = mat.Open(*outputFile, "w7.3"); err != nil {
			fail(err)
		}
		defer output.Close()
	}

	if err := command(config, input, output); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Printf("Error: %s.\n\n", err)

	fmt.Printf("Usage: %s [options]", os.Args[0])
	fmt.Printf(`

Options:
    -c <FILE.json>  - a configuration file (required)
    -i <FILE.mat>   - an input file
    -o <FILE.mat>   - an output file
`)

	os.Exit(1)
}
