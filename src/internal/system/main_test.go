package system

import (
	"testing"

	"github.com/ready-steady/assert"
	"github.com/turing-complete/laboratory/src/internal/config"
)

func TestNew(t *testing.T) {
	config, _ := config.New("fixtures/002_020.json")

	system, _ := New(&config.System)
	assert.Equal(system.Platform.Len(), 2, t)
	assert.Equal(system.Application.Len(), 20, t)

	schedule := system.schedule
	assert.Equal(schedule.Mapping, []uint{
		0, 1, 0, 0, 1, 1, 1, 0, 0, 1,
		1, 0, 0, 0, 0, 1, 1, 1, 1, 1,
	}, t)
	assert.Equal(schedule.Order, []uint{
		0, 1, 2, 9, 12, 16, 18, 14, 17, 13,
		15, 3, 5, 11, 19, 8, 7, 6, 4, 10,
	}, t)
	assert.EqualWithin(schedule.Start, []float64{
		0.000, 0.010, 0.013, 0.187, 0.265, 0.218, 0.262, 0.260, 0.242, 0.051,
		0.267, 0.237, 0.079, 0.152, 0.113, 0.170, 0.079, 0.141, 0.113, 0.242,
	}, 1e-15, t)
	assert.EqualWithin(schedule.Span, 0.291, 1e-15, t)
}
