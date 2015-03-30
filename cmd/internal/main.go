package internal

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/ready-steady/hdf5"
)

var (
	configFile  = flag.String("c", "", "a configuration file in JSON")
	inputFile   = flag.String("i", "", "a data file in HDF5 (typically an input)")
	outputFile  = flag.String("o", "", "a data file in HDF5 (typically an output)")
	profileFile = flag.String("p", "", "a file for dumping profiling information")
	verbose     = flag.Bool("v", false, "a flag for displaying diagnostic information")
)

func Run(command func(*Config, *hdf5.File, *hdf5.File) error) {
	flag.Usage = usage
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

	var input, output *hdf5.File

	if len(*configFile) == 0 {
		fail(errors.New("a configuration file is required"))
	}
	config, err := NewConfig(*configFile)
	if err != nil {
		fail(err)
	}
	if *verbose {
		config.Verbose = true
	}

	if len(*inputFile) > 0 {
		if input, err = hdf5.Open(*inputFile); err != nil {
			fail(err)
		}
		defer input.Close()
	}

	if len(*outputFile) > 0 {
		if _, err = os.Stat(*outputFile); os.IsNotExist(err) {
			output, err = hdf5.Create(*outputFile)
		} else {
			output, err = hdf5.Open(*outputFile)
		}
		if err != nil {
			fail(err)
		}
		defer output.Close()
	}

	if err = command(config, input, output); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Printf("Error: %s.\n", err)
	os.Exit(1)
}

func usage() {
	fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
	fmt.Printf("Options:\n")
	flag.PrintDefaults()
	os.Exit(1)
}
