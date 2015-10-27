package main

import (
	"runtime"
	"sync"
	"testing"

	"github.com/ready-steady/probability"
	"github.com/simulated-reality/laboratory/cmd/internal"
)

func BenchmarkInvokeJobQueue(b *testing.B) {
	benchmarkInvoke(invokeJobQueue, b)
}

func BenchmarkInvokeNoJobQueue(b *testing.B) {
	benchmarkInvoke(invokeNoJobQueue, b)
}

func benchmarkInvoke(invoke func(internal.Target, []float64) []float64, b *testing.B) {
	const (
		sampleCount = 1000
	)

	config, _ := internal.NewConfig("fixtures/002_020_profile.json")
	problem, _ := internal.NewProblem(config)
	target, _ := internal.NewTarget(problem)

	ni, _ := target.Dimensions()

	points := probability.Sample(probability.NewUniform(0, 1),
		probability.NewGenerator(0), sampleCount*ni)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		invoke(target, points)
	}
}

func invokeJobQueue(target internal.Target, points []float64) []float64 {
	return internal.Invoke(target, points, uint(runtime.GOMAXPROCS(0)))
}

func invokeNoJobQueue(target internal.Target, points []float64) []float64 {
	ic, oc := target.Dimensions()
	pc := uint(len(points)) / ic

	values := make([]float64, pc*oc)
	group := sync.WaitGroup{}
	group.Add(int(pc))

	for i := uint(0); i < pc; i++ {
		go func(point, value []float64) {
			target.Compute(point, value)
			group.Done()
		}(points[i*ic:(i+1)*ic], values[i*oc:(i+1)*oc])
	}

	group.Wait()

	return values
}
