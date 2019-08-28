package cachedrepo_test

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	cp "github.com/yeqown/cached-repository"
	"github.com/yeqown/cached-repository/lru"

	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	c cp.CacheAlgor
}

func (su *testSuite) SetupTest() {
	lruKC, err := lru.NewLRUK(2, 10, 20, func(k, v interface{}) {
		fmt.Printf("onEvcit: %v, %v", k, v)
	})
	if err != nil {
		panic(err)
	}
	su.c = cp.New(lruKC)
}

func (su *testSuite) TestNormalFunction() {
	su.c.Put("key1", 1)
	v, ok := su.c.Get("key1")
	su.Equal(false, ok)
	su.Equal(nil, v)

	su.c.Put("key1", 1)
	v, ok = su.c.Get("key1")
	su.Equal(true, ok)
	su.Equal(1, v)

	su.c.Put("key1", 2)
	v, ok = su.c.Get("key1")
	su.Equal(true, ok)
	su.Equal(2, v)

	for i := 2; i < 12; i++ {
		su.c.Put(fmt.Sprintf("key%d", i), i)
		su.c.Put(fmt.Sprintf("key%d", i), i)
	}
	v, ok = su.c.Get("key1")
	su.Equal(false, ok)
	su.Equal(nil, v)
}

func (su *testSuite) TestConcurrent() {
	timer := time.NewTimer(10 * time.Second)
	rand.Seed(time.Now().UnixNano())
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-timer.C:
				return
			default:
				v := rand.Int31()
				su.c.Put("key", v)
				log.Printf("put key=[%v], value=[%v]", "key", v)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// put into random data
	go func() {
		defer wg.Done()
		for {
			select {
			case <-timer.C:
				return
			default:
				v, ok := su.c.Get("key")
				log.Printf("get key, v=[%v], ok=[%v]", v, ok)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	wg.Wait()
	log.Println("done")
}

func Test_NormalCache(t *testing.T) {
	suite.Run(t, new(testSuite))
}
