package correlation

import (
	"math"

	"github.com/ready-steady/simulation/system"
)

func Compute(application *system.Application, index []uint, length float64) []float64 {
	nt, nd := uint(len(application.Tasks)), uint(len(index))

	distance := measure(application)
	C := make([]float64, nd*nd)

	for i := uint(0); i < nd; i++ {
		C[i*nd+i] = 1

		if length == 0 {
			continue
		}

		for j := i + 1; j < nd; j++ {
			d := distance[index[i]*nt+index[j]]
			C[j*nd+i] = math.Exp(-d * d / (length * length))
			C[i*nd+j] = C[j*nd+i]
		}
	}

	return C
}

func measure(application *system.Application) []float64 {
	nt := uint(len(application.Tasks))

	depth := explore(application)

	index := make([]uint, nt)
	count := make([]uint, nt)
	for i, d := range depth {
		index[i] = count[d]
		count[d]++
	}

	distance := make([]float64, nt*nt)

	for i := uint(0); i < nt; i++ {
		for j := i + 1; j < nt; j++ {
			xi := float64(index[i]) - float64(count[depth[i]])/2.0
			yi := float64(depth[i])

			xj := float64(index[j]) - float64(count[depth[j]])/2.0
			yj := float64(depth[j])

			distance[j*nt+i] = math.Sqrt((xi-xj)*(xi-xj) + (yi-yj)*(yi-yj))
			distance[i*nt+j] = distance[j*nt+i]
		}
	}

	return distance
}

func explore(application *system.Application) []uint {
	nt := uint(len(application.Tasks))
	depth := make([]uint, nt)

	for _, l := range application.Leafs() {
		ascend(application, depth, l)
	}

	return depth
}

func ascend(application *system.Application, depth []uint, f uint) {
	max := uint(0)

	for _, p := range application.Tasks[f].Parents {
		if depth[p] == 0 {
			ascend(application, depth, p)
		}
		if max < depth[p]+1 {
			max = depth[p] + 1
		}
	}

	depth[f] = max
}
