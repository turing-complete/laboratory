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
		var mode string
		if _, err = os.Stat(*outputFile); os.IsNotExist(err) {
			mode = "w7.3"
		} else {
			mode = "u"
		}
		if output, err = mat.Open(*outputFile, mode); err != nil {
			fail(err)
		}
		defer output.Close()
	}

	if err = command(config, input, output); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Printf("Error: %s.\n\n", err)

	fmt.Printf("Usage: %s [options]", os.Args[0])
	fmt.Printf(`

Options:
    -c <FILE.json>  - a configuration file (required)
    -i <FILE.mat>   - a data file, typically an input
    -o <FILE.mat>   - a date file, typically an output
`)

	os.Exit(1)
}
