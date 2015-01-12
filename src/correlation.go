package main

import (
	"math"

	"github.com/ready-steady/persim/system"
)

func correlate(app *system.Application, index []uint16, length float64) []float64 {
	tc, dc := uint16(len(app.Tasks)), uint16(len(index))

	distance := measure(app)
	C := make([]float64, dc*dc)

	for i := uint16(0); i < dc; i++ {
		C[i*dc+i] = 1
		for j := i + 1; j < dc; j++ {
			d := distance[index[i]*tc+index[j]]
			C[j*dc+i] = math.Exp(-d * d / (length * length))
			C[i*dc+j] = C[j*dc+i]
		}
	}

	return C
}

func measure(app *system.Application) []float64 {
	tc := uint16(len(app.Tasks))

	depth := explore(app)

	index := make([]uint16, tc)
	count := make([]uint16, tc)
	for i, d := range depth {
		index[i] = count[d]
		count[d]++
	}

	distance := make([]float64, tc*tc)

	for i := uint16(0); i < tc; i++ {
		for j := i + 1; j < tc; j++ {
			xi := float64(index[i]) - float64(count[depth[i]])/2.0
			yi := float64(depth[i])

			xj := float64(index[j]) - float64(count[depth[j]])/2.0
			yj := float64(depth[j])

			distance[j*tc+i] = math.Sqrt((xi-xj)*(xi-xj) + (yi-yj)*(yi-yj))
			distance[i*tc+j] = distance[j*tc+i]
		}
	}

	return distance
}

func explore(app *system.Application) []uint16 {
	tc := uint16(len(app.Tasks))
	depth := make([]uint16, tc)

	for _, l := range app.Leafs() {
		ascend(app, depth, l)
	}

	return depth
}

func ascend(app *system.Application, depth []uint16, f uint16) {
	max := uint16(0)

	for _, p := range app.Tasks[f].Parents {
		if depth[p] == 0 {
			ascend(app, depth, p)
		}
		if max < depth[p]+1 {
			max = depth[p] + 1
		}
	}

	depth[f] = max
}
