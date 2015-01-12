package main

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestNewProblemGeneral(t *testing.T) {
	config, err := loadConfig("fixtures/002_020.json")
	assert.Success(err, t)

	p, err := newProblem(config)
	assert.Success(err, t)

	assert.Equal(p.cc, uint32(2), t)
	assert.Equal(p.tc, uint32(20), t)

	assert.Equal(p.sched.Mapping, []uint16{
		0, 1, 0, 0, 1, 1, 1, 0, 0, 1,
		1, 0, 0, 0, 0, 1, 1, 1, 1, 1,
	}, t)
	assert.Equal(p.sched.Order, []uint16{
		0, 1, 2, 9, 12, 16, 18, 14, 17, 13,
		15, 3, 5, 11, 19, 8, 7, 6, 4, 10,
	}, t)
	assert.AlmostEqual(p.sched.Start, []float64{
		0.000, 0.010, 0.013, 0.187, 0.265, 0.218, 0.262, 0.260, 0.242, 0.051,
		0.267, 0.237, 0.079, 0.152, 0.113, 0.170, 0.079, 0.141, 0.113, 0.242,
	}, t)
	assert.AlmostEqual(p.sched.Span, 0.291, t)

	assert.Equal(p.sc, uint32(29100), t)
}

func TestNewProblemNoCores(t *testing.T) {
	config, _ := loadConfig("fixtures/002_020.json")
	config.CoreIndex = []uint16{}

	p, _ := newProblem(config)
	assert.Equal(p.config.CoreIndex, index(2), t)
}

func TestNewProblemNoTasks(t *testing.T) {
	config, _ := loadConfig("fixtures/002_020.json")
	config.TaskIndex = []uint16{}

	p, _ := newProblem(config)
	assert.Equal(p.config.TaskIndex, index(20), t)
}

func TestNewProblemProbModel(t *testing.T) {
	config, _ := loadConfig("fixtures/002_020.json")
	p, _ := newProblem(config)

	assert.Equal(p.uc, uint32(20), t)
	assert.Equal(p.zc, uint32(3), t)

	assert.Equal(len(p.trans), 3*20, t)

	delay := make([]float64, 20)
	for i := 0; i < 20; i++ {
		assert.Equal(p.marginals[i].InvCDF(0), 0.0, t)
		delay[i] = p.marginals[i].InvCDF(1)
	}
	assert.AlmostEqual(delay, []float64{
		0.0020, 0.0006, 0.0076, 0.0062, 0.0004, 0.0038, 0.0006, 0.0062, 0.0036, 0.0056,
		0.0038, 0.0010, 0.0068, 0.0070, 0.0078, 0.0044, 0.0004, 0.0058, 0.0056, 0.0040,
	}, t)

	assert.Equal(p.ic, uint32(3+1), t)
	assert.Equal(p.oc, uint32(1), t)
}

func BenchmarkSolve(b *testing.B) {
	c, _ := loadConfig("fixtures/002_020.json")

	for i := 0; i < b.N; i++ {
		p, _ := newProblem(c)
		p.solve()
	}
}
