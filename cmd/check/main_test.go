package main

import (
	"runtime"
	"testing"

	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"

	"../internal"
)

var (
	workerCount = uint32(runtime.GOMAXPROCS(0))
)

func BenchmarkInvokeWithoutChannels(b *testing.B) {
	benchmarkInvoke(invoke, b)
}

func BenchmarkInvokeWithChannels(b *testing.B) {
	benchmarkInvoke(invokeChannel, b)
}

func benchmarkInvoke(invoke func(internal.Target, []float64) []float64, b *testing.B) {
	const (
		sampleCount = 10000
	)

	problem, _ := internal.SetupProblem("fixtures/002_020.json")
	target, _ := internal.SetupTarget(problem)

	ic, _ := target.InputsOutputs()

	points := probability.Sample(uniform.New(0, 1), sampleCount*ic)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		invoke(target, points)
	}
}

func invokeChannel(target internal.Target, points []float64) []float64 {
	ic, oc := target.InputsOutputs()
	pc := uint32(len(points)) / ic

	values := make([]float64, pc*oc)
	jobs := make(chan uint32, pc)
	done := make(chan bool, pc)

	for i := uint32(0); i < workerCount; i++ {
		go func() {
			for j := range jobs {
				target.Evaluate(points[j*ic:(j+1)*ic], values[j*oc:(j+1)*oc], nil)
				done <- true
			}
		}()
	}

	for i := uint32(0); i < pc; i++ {
		jobs <- i
	}
	for i := uint32(0); i < pc; i++ {
		<-done
	}

	close(jobs)

	return values
}
