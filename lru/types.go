package lru

// EvictCallback .
type EvictCallback func(k, v interface{})

// IterFunc .
type IterFunc func(k, v interface{})

type entry struct {
	Key   interface{}
	Value interface{}
}

type historyEntry struct {
	Key     interface{}
	Value   interface{}
	Visited uint
}

// Cache is the interface for simple LRU cache.
type Cache interface {
	// Puts a value to the cache, returns true if an eviction occurred and
	// updates the "recently used"-ness of the key.
	Put(key, value interface{}) bool

	// Returns key's value from the cache and
	// updates the "recently used"-ness of the key. #value, isFound
	Get(key interface{}) (value interface{}, ok bool)

	// Removes a key from the cache.
	Remove(key interface{}) bool

	// Peeks a key
	// Returns key's value without updating the "recently used"-ness of the key.
	Peek(key interface{}) (value interface{}, ok bool)

	// Returns the oldest entry from the cache. #key, value, isFound
	Oldest() (interface{}, interface{}, bool)

	// Returns a slice of the keys in the cache, from oldest to newest.
	Keys() []interface{}

	// Returns the number of items in the cache.
	Len() int

	// iter all key and items in cache
	Iter(f IterFunc)

	// Clears all cache entries.
	Purge()
}
