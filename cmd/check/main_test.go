package main

import (
	"testing"

	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"

	"../internal"
)

func BenchmarkInvokeNoChannels(b *testing.B) {
	benchmarkInvoke(invokeNoChannel, b)
}

func BenchmarkInvokeWithChannels(b *testing.B) {
	benchmarkInvoke(invoke, b)
}

func benchmarkInvoke(invoke func(internal.Target, []float64) []float64, b *testing.B) {
	const (
		sampleCount = 10000
	)

	problem, _ := internal.NewProblem("fixtures/002_020.json")
	target, _ := internal.NewTarget(problem)

	ic, _ := target.InputsOutputs()

	points := probability.Sample(uniform.New(0, 1), sampleCount*ic)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		invoke(target, points)
	}
}

func invokeNoChannel(target internal.Target, points []float64) []float64 {
	ic, oc := target.InputsOutputs()
	pc := uint32(len(points)) / ic

	values := make([]float64, pc*oc)
	done := make(chan bool, pc)

	for i := uint32(0); i < pc; i++ {
		go func(point, value []float64) {
			target.Evaluate(point, value, nil)
			done <- true
		}(points[i*ic:(i+1)*ic], values[i*oc:(i+1)*oc])
	}
	for i := uint32(0); i < pc; i++ {
		<-done
	}

	return values
}