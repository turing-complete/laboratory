package cache

import (
	"fmt"
	"unsafe"
)

// No cuncurrent access for now!
type Cache struct {
	hc uint32
	mc uint32

	length uint32
	buffer []byte

	storage map[string][]float64
}

func (c *Cache) String() string {
	return fmt.Sprintf("Cache{hits: %d (%.2f%%), misses: %d (%.2f%%)}",
		c.hc, float64(c.hc)/float64(c.hc+c.mc)*100,
		c.mc, float64(c.mc)/float64(c.hc+c.mc)*100)
}

func New(length uint32, space uint32) *Cache {
	return &Cache{
		length:  length,
		buffer:  make([]byte, 8*length),
		storage: make(map[string][]float64, space),
	}
}

func (c *Cache) Key(sequence []uint64) string {
	for i := uint32(0); i < c.length; i++ {
		*(*uint64)(unsafe.Pointer(&c.buffer[8*i])) = sequence[i]
	}
	return string(c.buffer)
}

func (c *Cache) Get(key string) []float64 {
	value := c.storage[key]

	if value != nil {
		c.hc++
	} else {
		c.mc++
	}

	return value
}

func (c *Cache) Set(key string, value []float64) {
	c.storage[key] = value
}

func (c *Cache) Flush() {
	c.storage = make(map[string][]float64, len(c.storage))
}
