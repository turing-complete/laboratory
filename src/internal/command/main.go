package command

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/target"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
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

func Setup(config *config.Config) (*system.System, *uncertainty.Uncertainty,
	target.Target, error) {

	system, err := system.New(&config.System)
	if err != nil {
		return nil, nil, nil, err
	}

	uncertainty, err := uncertainty.New(system, &config.Uncertainty)
	if err != nil {
		return nil, nil, nil, err
	}

	target, err := target.New(system, uncertainty, &config.Target)
	if err != nil {
		return nil, nil, nil, err
	}

	return system, uncertainty, target, nil
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
