package cache

import (
	"fmt"
	"reflect"
	"unsafe"
)

// No cuncurrent access for now!
type Cache struct {
	depth   int
	mapping map[string][]float64

	hc uint32
	mc uint32
}

func (c *Cache) String() string {
	return fmt.Sprintf("Cache{hits: %d (%.2f%%), misses: %d (%.2f%%)}",
		c.hc, float64(c.hc)/float64(c.hc+c.mc)*100,
		c.mc, float64(c.mc)/float64(c.hc+c.mc)*100)
}

func New(depth uint32, capacity uint32) *Cache {
	return &Cache{
		depth:   int(depth),
		mapping: make(map[string][]float64, capacity),
	}
}

func (c *Cache) Key(trace []uint64) string {
	const (
		sizeOfUInt64 = 8
	)

	sliceHeader := *(*reflect.SliceHeader)(unsafe.Pointer(&trace))

	stringHeader := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sizeOfUInt64 * c.depth,
	}

	return *(*string)(unsafe.Pointer(&stringHeader))
}

func (c *Cache) Get(key string) []float64 {
	value := c.mapping[key]

	if value != nil {
		c.hc++
	} else {
		c.mc++
	}

	return value
}

func (c *Cache) Set(key string, value []float64) {
	c.mapping[key] = value
}

func (c *Cache) Flush() {
	c.mapping = make(map[string][]float64, len(c.mapping))
}
