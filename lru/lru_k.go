package lru

import (
	"container/list"
	"errors"
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
