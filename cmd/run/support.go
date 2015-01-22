package main

import (
	"time"
)

func track(work func()) time.Duration {
	start := time.Now()
	work()
	return time.Now().Sub(start)
}
