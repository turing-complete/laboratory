package main

import (
	"math"
	"testing"

	"github.com/ready-steady/persim/system"
	"github.com/ready-steady/stats/decomp"
	"github.com/ready-steady/support/assert"
)

func TestCorrelate(t *testing.T) {
	_, app, _ := system.Load("fixtures/002_020.tgff")

	C := correlate(app, index(20), 2)
	_, _, err := decomp.CovPCA(C, 20)
	assert.Success(err, t)

	C = correlate(app, index(1), 2)
	assert.Equal(C, []float64{1}, t)
}

func TestMeasure(t *testing.T) {
	_, app, _ := system.Load("fixtures/002_020.tgff")
	distance := measure(app)

	cases := []struct {
		i uint16
		j uint16
		d float64
	}{
		{0, 1, 1},
		{0, 7, 3},
		{0, 18, math.Sqrt(5*5 + 0.5*0.5)},
		{1, 2, math.Sqrt(1*1 + 1*1)},
		{1, 3, 1},
		{2, 3, 1},
		{3, 9, math.Sqrt(1*1 + 2*2)},
		{8, 9, 1},
	}

	for _, c := range cases {
		assert.Equal(distance[20*c.i+c.j], c.d, t)
	}
}

func TestExplore(t *testing.T) {
	_, app, _ := system.Load("fixtures/002_020.tgff")
	depth := explore(app)

	assert.Equal(depth, []uint16{
		0,
		1,
		2, 2, 2,
		3, 3, 3, 3, 3,
		4, 4, 4, 4, 4, 4, 4, 4,
		5, 5,
	}, t)
}

func BenchmarkCorrelate(b *testing.B) {
	_, app, _ := system.Load("fixtures/002_020.tgff")
	index := index(20)

	for i := 0; i < b.N; i++ {
		correlate(app, index, 2)
	}
}
