package lru

import (
	"container/list"
)

var (
	_ Cache = &LRU{}
)

// LRU .
// TODO: goroutine safe
type LRU struct {
	size        uint                          // max size
	cache       *list.List                    // doubly linked list
	itemRecords map[interface{}]*list.Element // item map, get faster
	onEvict     EvictCallback                 // callback func
}

// NewLRU constructs an LRU of the given size
func NewLRU(size uint, onEvict EvictCallback) (*LRU, error) {
	c := &LRU{
		size:        size,
		cache:       list.New(),
		itemRecords: make(map[interface{}]*list.Element),
		onEvict:     onEvict,
	}
	return c, nil
}

// Purge is used to completely clear the cache.
func (c *LRU) Purge() {
	for k, v := range c.itemRecords {
		if c.onEvict != nil {
			c.onEvict(k, v.Value.(*entry).Value)
		}
		delete(c.itemRecords, k)
	}
	c.cache.Init()
}

// Put adds a value to the cache.  Returns true if an eviction occurred.
func (c *LRU) Put(key, value interface{}) (evicted bool) {
	// Check for existing item
	if item, ok := c.itemRecords[key]; ok {
		c.cache.MoveToFront(item)
		item.Value.(*entry).Value = value
		return false
	}

	// Add new item
	ent := &entry{key, value}
	item := c.cache.PushFront(ent)
	c.itemRecords[key] = item

	// Verify size not exceeded
	if evicted = c.cache.Len() > int(c.size); evicted {
		c.removeOldest()
	}

	return evicted
}

// Get looks up a key's value from the cache.
func (c *LRU) Get(key interface{}) (value interface{}, ok bool) {
	if item, ok := c.itemRecords[key]; ok {
		c.cache.MoveToFront(item)
		// if item.Value.(*entry) == nil {
		// 	return nil, false
		// }
		return item.Value.(*entry).Value, true
	}
	return
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *LRU) Peek(key interface{}) (value interface{}, ok bool) {
	var item *list.Element
	if item, ok = c.itemRecords[key]; ok {
		return item.Value.(*entry).Value, true
	}
	return nil, ok
}

// Remove removes the provided key from the cache, returning if the
// key was contained.
func (c *LRU) Remove(key interface{}) (present bool) {
	if item, ok := c.itemRecords[key]; ok {
		c.removeElement(item)
		return true
	}
	return false
}

// RemoveOldest removes the oldest item from the cache.
func (c *LRU) RemoveOldest() (key interface{}, value interface{}, ok bool) {
	item := c.cache.Back()
	if item != nil {
		c.removeElement(item)
		ent := item.Value.(*entry)
		return ent.Key, ent.Value, true
	}
	return nil, nil, false
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *LRU) Keys() []interface{} {
	keys := make([]interface{}, len(c.itemRecords))
	i := 0
	for item := c.cache.Back(); item != nil; item = item.Prev() {
		keys[i] = item.Value.(*entry).Key
		i++
	}
	return keys
}

// Len returns the number of itemRecords in the cache.
func (c *LRU) Len() int {
	return c.cache.Len()
}

// Oldest returns the oldest item in the cache.
func (c *LRU) Oldest() (interface{}, interface{}, bool) {
	if c.cache.Len() == 0 {
		return nil, nil, false
	}

	item := c.cache.Back()
	ent := item.Value.(*entry)
	return ent.Value, ent.Value, true
}

// Iter .
func (c *LRU) Iter(f IterFunc) {
	for item := c.cache.Back(); item != nil; item = item.Prev() {
		ent := item.Value.(*entry)
		f(ent.Key, ent.Value)
	}
}

// removeOldest removes the oldest item from the cache.
func (c *LRU) removeOldest() {
	item := c.cache.Back()
	if item != nil {
		c.removeElement(item)
	}
}

// removeElement is used to remove a given list element from the cache
func (c *LRU) removeElement(item *list.Element) {
	c.cache.Remove(item)
	ent := item.Value.(*entry)
	delete(c.itemRecords, ent.Key)
	if c.onEvict != nil {
		c.onEvict(ent.Key, ent.Value)
	}
}
