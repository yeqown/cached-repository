package main

import (
	"fmt"

	cachedrepo "github.com/yeqown/cached-repository"
	"github.com/yeqown/cached-repository/lru"
)

func main() {
	c, err := lru.NewLRUK(2, 2, 10, nil)
	if err != nil {
		panic(err)
	}
	ca := cachedrepo.New(c)
	v, ok := ca.Get("key1")
	fmt.Printf("ca.Get('key1')=%v,%v\n", v, ok)

	// put key1, twice means 2 times visit
	ca.Put("key1", "value1")
	ca.Put("key1", "value1")
	v, ok = ca.Get("key1")
	fmt.Printf("ca.Get('key1')=%v,%v\n", v, ok)

	// update key1
	ca.Update("key1", "value111")
	v, ok = ca.Get("key1")
	fmt.Printf("ca.Get('key1')=%v,%v\n", v, ok)

	// put key2, key3
	ca.Put("key2", "value2")
	ca.Put("key2", "value2")
	ca.Put("key3", "value3")
	ca.Put("key3", "value3")

	// query key1 again
	v, ok = ca.Get("key1")
	fmt.Printf("ca.Get('key1')=%v,%v\n", v, ok)

	// query key2
	v, ok = ca.Get("key2")
	fmt.Printf("ca.Get('key2')=%v,%v\n", v, ok)

	// query key3
	v, ok = ca.Get("key3")
	fmt.Printf("ca.Get('key3')=%v,%v\n", v, ok)
}
