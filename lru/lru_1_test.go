package lru_test

import (
	"fmt"
	"testing"

	"github.com/yeqown/cached-repository/lru"
)

func Test_LRU1(t *testing.T) {
	cache, err := lru.NewLRU(1, func(k, v interface{}) {
		fmt.Printf("onEvict: k: %v, v: %v\n", k, v)
	})

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var (
		hit bool
	)

	cache.Put(2, 1)
	if _, hit = cache.Get(2); !hit {
		t.Error("should not hit 2")
		t.FailNow()
	}

	// replacing
	cache.Put(3, 2)
	if _, hit = cache.Get(2); hit {
		t.Error("should not hit 2")
		t.FailNow()
	}
	if _, hit = cache.Get(3); !hit {
		t.Error("should hit 3")
		t.FailNow()
	}

	// update
	cache.Put(3, 1)
	if v, hit := cache.Get(3); !hit || v.(int) != 1 {
		t.Error("should hit 3 and v = 1")
		t.FailNow()
	}
}
