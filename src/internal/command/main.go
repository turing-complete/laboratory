package command

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/turing-complete/laboratory/src/internal/config"
)

var (
	configFile  = flag.String("c", "", "a configuration file (required)")
	profileFile = flag.String("p", "", "an output file for profiling information")
	verbose     = flag.Bool("v", false, "a flag for displaying diagnostic information")
)

type null struct{}

func (_ null) Write(buffer []byte) (int, error) {
	return len(buffer), nil
}

func Run(function func(*config.Config) error) {
	flag.Usage = usage
	flag.Parse()

	if len(*profileFile) > 0 {
		profile, err := os.Create(*profileFile)
		if err != nil {
			fail(errors.New("cannot enable profiling"))
		}
		pprof.StartCPUProfile(profile)
		defer pprof.StopCPUProfile()
	}

	if len(*configFile) == 0 {
		fail(errors.New("expected a filename"))
	}
	config, err := config.New(*configFile)
	if err != nil {
		fail(err)
	}
	if *verbose {
		config.Verbose = true
	}
	if !config.Verbose {
		log.SetOutput(null{})
	}

	if err = function(config); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Errorf("Error: %s.\n", err)
	os.Exit(1)
}

func usage() {
	fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
	fmt.Printf("Options:\n")
	flag.PrintDefaults()
	os.Exit(1)
}
