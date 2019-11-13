# cached-repository
[![Go Report Card](https://goreportcard.com/badge/github.com/yeqown/cached-repository)](https://goreportcard.com/report/github.com/yeqown/cached-repository) [![](https://godoc.org/github.com/yeqown/cached-repository?status.svg)](https://godoc.org/github.com/yeqown/cached-repository) ![](https://img.shields.io/badge/LICENSE-MIT-blue.svg)

a basic repository to support cached data in memory and is based LRU-K cache replacing algorithm

### TODOs

* [x] LRU-1 & LRU-K

* [x] Cached-Repository [demo](./examples/custom-cache-manage/main.go)

* [x] `LRU-K` concurrent safe

* [ ] `LRU-1` concurrent safe

### Quick Start

`simple`
```go
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

	// put key2, key3 twice
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
```

```sh
ca.Get('key1')=<nil>,false
ca.Get('key1')=value1,true
ca.Get('key1')=value111,true
ca.Get('key1')=<nil>,false
ca.Get('key2')=value2,true
ca.Get('key3')=value3,true
```


### Examples

* [simple](/example/simple)
* [cache repository](/examples/custome-cache-manage)
