package internal

import (
	"time"
)

func Track(work func()) time.Duration {
	start := time.Now()
	work()
	return time.Now().Sub(start)
}
