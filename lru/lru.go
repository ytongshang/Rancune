package lru

import (
	"container/list"
)

type Key interface{}

type SizeOf func(key Key, value interface{}) int64

var DefaultSizeOf = func(key Key, value interface{}) int64 {
	return 1
}

type EntryRemoved func(evicted bool, key Key, oldValue interface{}, newValue interface{})

type Cache struct {
	MaxSize int64
	size    int64

	ll    *list.List
	cache map[interface{}]*list.Element

	EntryRemoved EntryRemoved

	SizeOf SizeOf
}

type entry struct {
	key   Key
	value interface{}
}

func NewLruCache(maxSize int64, entryremoved EntryRemoved, sizeof SizeOf) *Cache {
	if maxSize < 0 {
		panic("lrucache, maxsize < 0")
	}
	if sizeof == nil {
		sizeof = DefaultSizeOf
	}
	return &Cache{
		MaxSize:      maxSize,
		size:         0,
		EntryRemoved: entryremoved,
		SizeOf:       sizeof,
		ll:           list.New(),
		cache:        make(map[interface{}]*list.Element),
	}
}

func (c *Cache) Put(key Key, value interface{}) interface{} {
	if key == nil || value == nil {
		panic("key == nil || value == nil")
	}
	c.lazyInit()
	c.size += c.safeSizeOf(key, value)
	var previous interface{}
	// 原来存在对应的缓存
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		element := ee.Value.(*entry)
		previous = element.value
		c.size -= c.safeSizeOf(key, previous)
		element.value = value
		c.entryRemoved(false, key, previous, value)
	} else {
		// 没有
		ele := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
	}
	c.trimToSize(c.MaxSize)
	return previous
}

func (c *Cache) Get(key Key) (interface{}, bool) {
	if key == nil {
		panic("key == nil")
	}
	if c.cache == nil {
		return nil, false
	}
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return nil, false
}

func (c *Cache) Remove(key Key) interface{} {
	if key == nil {
		panic("key == nil")
	}
	if c.cache == nil {
		return nil
	}
	if ele, ok := c.cache[key]; ok {
		return c.removeElement(ele)
	}
	return nil
}

func (c *Cache) RemoveEldest() interface{} {
	if c.cache == nil {
		return nil
	}
	ele := c.ll.Back()
	if ele != nil {
		return c.removeElement(ele)
	}
	return nil
}

func (c *Cache) Clear() {
	if c.EntryRemoved != nil {
		for _, ele := range c.cache {
			entry := ele.Value.(*entry)
			c.entryRemoved(false, entry.key, entry.value, nil)
		}
	}
	c.ll = nil
	c.cache = nil
	c.size = 0
}

func (c *Cache) Size() int64 {
	return c.size
}

func (c *Cache) removeElement(e *list.Element) interface{} {
	entry := e.Value.(*entry)
	delete(c.cache, entry.key)
	c.ll.Remove(e)
	c.size -= c.safeSizeOf(entry.key, entry.value)
	c.entryRemoved(false, entry.key, entry.value, nil)
	return entry.value
}

func (c *Cache) lazyInit() {
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}
}

func (c *Cache) safeSizeOf(key Key, value interface{}) int64 {
	if c.SizeOf == nil {
		c.SizeOf = DefaultSizeOf
	}
	size := c.SizeOf(key, value)
	return size
}

func (c *Cache) entryRemoved(evicted bool, key Key, oldValue interface{}, newValue interface{}) {
	if c.EntryRemoved != nil {
		c.EntryRemoved(evicted, key, oldValue, newValue)
	}
}

func (c *Cache) trimToSize(maxSize int64) {
	if c.cache == nil {
		return
	}
	for {
		if c.size < maxSize {
			break
		}
		eldest := c.ll.Back()
		if eldest == nil {
			break
		}
		c.removeElement(eldest)
	}
}
