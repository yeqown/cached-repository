package lru

import (
	"container/list"
	"errors"
)

var (
	_ Cache = K{}
)

// K . means lru-k
type K struct {
	K       uint          // the K setting
	onEvict EvictCallback // evict callback

	hSize        uint                          // historyMax - used = historyRest
	history      *list.List                    // history doubly linked list
	historyItems map[interface{}]*list.Element // history get op O(1)

	size       uint                          // max - used = rest
	cache      *list.List                    // cache doubly linked list, save
	cacheItems map[interface{}]*list.Element // cache get op O(1)
}

// NewK .
func NewK(k, size, hSize uint, onEvict EvictCallback) (*K, error) {

	if k < 2 {
		return nil, errors.New("k is suggested bigger than 1, otherwise using LRU")
	}

	if hSize < size {
		hSize = size * ((size % 3) + 1) // why would i set this?
	}

	return &K{
		K:            k,
		onEvict:      onEvict,
		hSize:        hSize,
		history:      list.New(),
		historyItems: make(map[interface{}]*list.Element),
		size:         size,
		cache:        list.New(),
		cacheItems:   make(map[interface{}]*list.Element),
	}, nil
}

// Put of K cache add or update
func (c K) Put(key, value interface{}) (evicted bool) {
	if item, ok := c.cacheItems[key]; ok {
		item.Value = value
		c.cache.MoveToFront(item)
		return
	}

	// not hit in cache, then add to history
	var hEnt = new(historyEntry)
	item, ok := c.historyItems[key]
	if ok {
		hEnt = item.Value.(*historyEntry)
		hEnt.Visited++
		item.Value = hEnt
		if hEnt.Visited >= c.K {
			// true: move from history into cache
			c.removeHistoryElement(item)
			return c.addElement(&entry{Key: key, Value: value})
		}
	} else {
		// true: not exists
		hEnt.Key = key
		hEnt.Value = value
		hEnt.Visited = 1
		item = c.addHistoryElement(hEnt)
	}
	// update history
	// c.historyItems[key] = item

	return false
}

// Get of K cache
func (c K) Get(key interface{}) (value interface{}, ok bool) {
	if item, ok := c.cacheItems[key]; ok {
		c.cache.MoveToFront(item)
		return item.Value.(*entry).Value, true
	}
	return nil, false
}

// Remove of K cache
func (c K) Remove(key interface{}) bool {
	if item, ok := c.cacheItems[key]; ok {
		c.removeElement(item)
		return true
	}
	return false
}

// Peek of K cache
func (c K) Peek(key interface{}) (value interface{}, ok bool) {
	var item *list.Element
	if item, ok = c.cacheItems[key]; ok {
		return item.Value.(*entry).Value, true
	}
	return nil, ok
}

// Oldest of K cache
func (c K) Oldest() (key, value interface{}, ok bool) {
	if c.cache == nil || c.cache.Len() == 0 {
		return nil, nil, false
	}

	item := c.cache.Back()
	ent := item.Value.(*entry)
	return ent.Value, ent.Value, true
}

// Keys of K cache
func (c K) Keys() []interface{} {
	keys := make([]interface{}, len(c.cacheItems))
	i := 0
	for item := c.cache.Back(); item != nil; item = item.Prev() {
		keys[i] = item.Value.(*entry).Key
		i++
	}
	return keys
}

// Len of K cache
func (c K) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.cache.Len()
}

// Iter of K cache
func (c K) Iter(f IterFunc) {
	for item := c.cache.Back(); item != nil; item = item.Prev() {
		ent := item.Value.(*entry)
		f(ent.Key, ent.Value)
	}
}

// Purge of K cache
func (c K) Purge() {
	for k, v := range c.cacheItems {
		if c.onEvict != nil {
			c.onEvict(k, v.Value.(*entry).Value)
		}
		delete(c.cacheItems, k)
	}
	c.cache.Init()

	for k := range c.historyItems {
		delete(c.historyItems, k)
	}
	c.history.Init()
}

func (c K) removeHistoryElement(item *list.Element) {
	c.hSize++
	ent := item.Value.(*historyEntry)
	c.history.Remove(item)
	delete(c.historyItems, ent.Key)
}

func (c K) addHistoryElement(hEnt *historyEntry) *list.Element {
	if c.size == 0 {
		c.removeHistoryElement(c.history.Back())
	}

	c.historyItems[hEnt.Key] = c.history.PushFront(hEnt)
	return c.historyItems[hEnt.Key]
}

func (c K) removeElement(item *list.Element) {
	c.size++
	ent := item.Value.(*entry)
	c.cache.Remove(item)
	delete(c.cacheItems, ent.Key)
	if c.onEvict != nil {
		c.onEvict(ent.Key, ent.Value)
	}
}

func (c K) addElement(ent *entry) (evicted bool) {
	if c.size == 0 {
		evicted = true
		c.removeElement(c.cache.Back())
	}

	c.cacheItems[ent.Key] = c.cache.PushFront(ent)
	return
}
