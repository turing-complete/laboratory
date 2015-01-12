package main

// #include <string.h>
import "C"

import (
	"unsafe"

	"github.com/ready-steady/linal/matrix"
)

type job struct {
	key   string
	data  []float64
	node  []float64
	value []float64
	done  chan<- result
}

type result struct {
	key  string
	data []float64
}

func serveEndToEndDelay(p *problem, jobs <-chan job) {
	uc, zc := p.uc, p.zc

	g, m := p.gaussian, p.marginals

	z := make([]float64, zc)
	u := make([]float64, uc)
	d := make([]float64, p.tc)

	for job := range jobs {
		span := job.data

		if span == nil {
			span = make([]float64, 1)

			// Independent uniform to independent Gaussian
			for i := uint32(0); i < zc; i++ {
				z[i] = g.InvCDF(job.node[i])
			}

			// Independent Gaussian to dependent Gaussian
			matrix.Multiply(p.trans, z, u, uc, zc, 1)

			// Dependent Gaussian to dependent uniform to dependent target
			for i, tid := range p.config.TaskIndex {
				d[tid] = m[i].InvCDF(g.CDF(u[i]))
			}

			span[0] = p.time.Recompute(p.sched, d).Span
		}

		job.done <- result{job.key, span}
	}
}

func serveTemperatureProfile(p *problem, jobs <-chan job) {
	cc, sc, uc, zc, oc := p.cc, p.sc, p.uc, p.zc, p.oc
	coreIndex := p.config.CoreIndex

	g, m := p.gaussian, p.marginals

	P := make([]float64, cc*sc)
	S := make([]float64, p.tempan.Nodes*sc)

	z := make([]float64, zc)
	u := make([]float64, uc)
	d := make([]float64, p.tc)

	for job := range jobs {
		Q := job.data

		if Q == nil {
			Q = make([]float64, cc*sc)

			// Independent uniform to independent Gaussian
			for i := uint32(0); i < zc; i++ {
				// NOTE: +1 for time
				z[i] = g.InvCDF(job.node[1+i])
			}

			// Independent Gaussian to dependent Gaussian
			matrix.Multiply(p.trans, z, u, uc, zc, 1)

			// Dependent Gaussian to dependent uniform to dependent target
			for i, tid := range p.config.TaskIndex {
				d[tid] = m[i].InvCDF(g.CDF(u[i]))
			}

			// FIXME: Bad, bad, bad!
			C.memset(unsafe.Pointer(&P[0]), 0, C.size_t(8*cc*sc))

			p.power.Compute(p.time.Recompute(p.sched, d), P, sc)
			p.tempan.ComputeTransient(P, Q, S, sc)
		}

		sid := uint32(job.node[0] * float64(sc-1))
		for i := uint32(0); i < oc; i++ {
			job.value[i] = Q[sid*cc+uint32(coreIndex[i])]
		}

		job.done <- result{job.key, Q}
	}
}
