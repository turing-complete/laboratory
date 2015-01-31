package pool

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

type Foo struct {
	bar int
	baz int
	qux int
}

func TestPutGet(t *testing.T) {
	pool := New(1, func() interface{} {
		return new(Foo)
	})

	pool.Put(&Foo{1, 2, 3})
	pool.Put(&Foo{4, 5, 6})

	assert.Equal(pool.Get().(*Foo), &Foo{1, 2, 3}, t)
	assert.Equal(pool.Get().(*Foo), &Foo{0, 0, 0}, t)
}
