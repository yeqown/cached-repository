package cachedrepo

// CacheAlgor is an interface implements different alg.
type CacheAlgor interface {
	GetByID()
	Query()
}

// LRUKCache .
type LRUKCache struct{}
