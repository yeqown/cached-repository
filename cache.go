package cachedrepo

import (
	"github.com/yeqown/cached-repository/lru"
)

// CacheAlgor is an interface implements different alg.
type CacheAlgor interface {
	Put(key, value interface{})
	Get(key interface{}) (value interface{}, ok bool)
	Update(key, value interface{})
	Delete(key interface{})
}

var (
	_ CacheAlgor = LRUCacheAlgor{}
)

// New .
func New(c lru.Cache) CacheAlgor {
	return LRUCacheAlgor{
		c: c,
	}
}

// LRUCacheAlgor .
type LRUCacheAlgor struct {
	c lru.Cache
}

// Put of LRUCacheAlgor
func (a LRUCacheAlgor) Put(key, value interface{}) {
	a.c.Put(key, value)
}

// Get of LRUCacheAlgor
func (a LRUCacheAlgor) Get(key interface{}) (value interface{}, ok bool) {
	return a.c.Get(key)
	// return nil, false
}

// Update of LRUCacheAlgor
func (a LRUCacheAlgor) Update(key, value interface{}) {
	a.c.Put(key, value)
}

// Delete of LRUCacheAlgor
func (a LRUCacheAlgor) Delete(key interface{}) {
	a.c.Remove(key)
}
