package cache

import (
	"math/rand"
	"testing"
)

func BenchmarkKey(b *testing.B) {
	const (
		count    = 10000
		depth    = 100
		capacity = 100
	)

	cache := New(depth, capacity)

	generator := rand.New(rand.NewSource(0))
	traces := make([]uint64, count*depth)
	for i := range traces {
		traces[i] = uint64(generator.Int63())
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < count; j++ {
			cache.Key(traces[j*depth : (j+1)*depth])
		}
	}
}
