package main

import (
	"testing"

	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"

	"../internal"
)

func BenchmarkInvokeJobQueue(b *testing.B) {
	benchmarkInvoke(invoke, b)
}

func BenchmarkInvokeNoJobQueue(b *testing.B) {
	benchmarkInvoke(invokeNoJobQueue, b)
}

func benchmarkInvoke(invoke func(internal.Target, []float64) []float64, b *testing.B) {
	const (
		sampleCount = 10000
	)

	config, _ := internal.NewConfig("fixtures/002_020.json")
	problem, _ := internal.NewProblem(config)
	target, _ := internal.NewTarget(problem)

	points := probability.Sample(uniform.New(0, 1), sampleCount*target.Inputs())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		invoke(target, points)
	}
}

func invokeNoJobQueue(target internal.Target, points []float64) []float64 {
	ic, oc := target.Inputs(), target.Outputs()
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
