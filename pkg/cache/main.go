package cache

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

type Cache struct {
	depth   int
	mapping map[string][]float64

	hc uint32
	mc uint32

	sync.Mutex
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
	c.Lock()

	value, ok := c.mapping[key]
	if ok {
		c.hc++
	} else {
		c.mc++
	}

	c.Unlock()

	return value
}

func (c *Cache) Set(key string, value []float64) {
	c.Lock()
	c.mapping[key] = value
	c.Unlock()
}
