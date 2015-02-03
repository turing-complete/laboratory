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
	defer c.mutex.Unlock()
	value, ok = c.mapping[key]
	return
}

func (c *Cache) Add(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.mapping[key] = value
}
