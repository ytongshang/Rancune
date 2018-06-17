package util

import "sync"

type Cache struct {
	sync.Mutex
	data map[string]string
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]string)}
}

func (cache *Cache) Get(key string) (value string, ok bool) {
	cache.Lock()
	defer cache.Unlock()
	value, ok = cache.data[key]
	return
}

func (cache *Cache) Put(key string, value string) {
	cache.Lock()
	defer cache.Unlock()
	cache.data[key] = value
}
