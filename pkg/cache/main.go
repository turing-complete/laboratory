package cache

import (
	"sync"
)

type Cache struct {
	mutex   sync.Mutex
	mapping map[string]interface{}
}

func New(capacity uint32) *Cache {
	return &Cache{
		mapping: make(map[string]interface{}, capacity),
	}
}

func (c *Cache) Get(key string) (value interface{}, ok bool) {
	c.mutex.Lock()
	value, ok = c.mapping[key]
	c.mutex.Unlock()
	return
}

func (c *Cache) Add(key string, value interface{}) {
	c.mutex.Lock()
	c.mapping[key] = value
	c.mutex.Unlock()
}

func (c *Cache) Flush() {
	c.mutex.Lock()
	c.mapping = make(map[string]interface{}, len(c.mapping))
	c.mutex.Unlock()
}
