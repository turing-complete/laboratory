package main

import (
	"fmt"
	"time"
)

func track(description string, verbose bool, work func()) {
	if verbose {
		fmt.Println(description)
	}

	start := time.Now()
	work()
	duration := time.Now().Sub(start)

	if verbose {
		fmt.Printf("Done in %v.\n", duration)
	}
}
