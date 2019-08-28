# cached-repository
a basic repository to support cached data in memory and is based LRU-K cache replacing algorithm

### TODOs

* [x] LRU-1 & LRU-K

* [x] Cached-Repository [demo](./examples/custom-cache-manage/main.go)

* [x] `LRU-K` concurrent safe

* [ ] `LRU-1` concurrent safe

### Demo

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
	"sync"

	cp "github.com/yeqown/cached-repository"
	"github.com/yeqown/infrastructure/framework/gormic"
	"github.com/yeqown/infrastructure/types"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/yeqown/cached-repository/lru"
)

func main() {
	repo, err := prepareData()
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		go func(){
		wg.Add(1)
			repo.Create(&userModel{
			Model: gorm.Model{
				ID: uint(i + 1),
			},
			Name:     fmt.Sprintf("name-%d", i+1),
			Province: fmt.Sprintf("province-%d", i+1),
			City:     fmt.Sprintf("city-%d", i+1),
		})
		wg.Done()
		}()
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000; i++ {
		go func() {
			wg.Add(1)
			id := uint(rand.Intn(10))
			if id == 0 {
				continue
			}
	
			v, err := repo.GetByID(id)
			if err != nil {
				fmt.Printf("err: %d , %v\n", id, err)
				continue
			}
	
			if v.ID != id ||
				v.Name != fmt.Sprintf("name-%d", id) ||
				v.Province != fmt.Sprintf("province-%d", id) ||
				v.City != fmt.Sprintf("city-%d", id) {
				fmt.Printf("err: not matched target with id[%d]: %v\n", v.ID, v)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func prepareData() (*MysqlRepo, error) {
	db, err := gormic.ConnectSqlite3(&types.SQLite3Config{
		Name: "./testdata/sqlite3.db",
	})
	if err != nil {
		return nil, err
	}

	db.DropTableIfExists(&userModel{})
	db.AutoMigrate(&userModel{})

	return NewMysqlRepo(db)
}

type userModel struct {
	gorm.Model
	Name     string `gorm:"column:name"`
	Province string `gorm:"column:province"`
	City     string `gorm:"column:city"`
}

// MysqlRepo .
type MysqlRepo struct {
	db   *gorm.DB
	calg cp.CacheAlgor
	// *cp.EmbedRepo
}

// NewMysqlRepo .
func NewMysqlRepo(db *gorm.DB) (*MysqlRepo, error) {
	c, err := lru.NewLRUK(2, 10, 20, func(k, v interface{}) {
		fmt.Printf("key: %v, value: %v\n", k, v)
	})
	if err != nil {
		return nil, err
	}

	return &MysqlRepo{
		db:   db,
		calg: cp.New(c),
	}, nil
}

// Create .
func (repo MysqlRepo) Create(m *userModel) error {
	return repo.db.Create(m).Error
}

// GetByID .
func (repo MysqlRepo) GetByID(id uint) (*userModel, error) {
	start := time.Now()
	defer func() {
		fmt.Printf("this queryid=%d cost: %d ns\n",id, time.Now().Sub(start).Nanoseconds())
	}()

	v, ok := repo.calg.Get(id)
	if ok {
		return v.(*userModel), nil
	}

	// actual find in DB
	m := new(userModel)
	if err := repo.db.Where("id = ?", id).First(m).Error; err != nil {
		return nil, err
	}

	repo.calg.Put(id, m)
	return m, nil
}

// Update .
func (repo MysqlRepo) Update(id uint, m *userModel) error {
	if err := repo.db.Where("id = ?", id).Update(m).Error; err != nil {
		return err
	}

	fmt.Printf("before: %v\n", m)
	m.ID = id
	if err := repo.db.First(m); err != nil {

	}
	fmt.Printf("after: %v\n", m)

	// update cache, ifcache hit id
	repo.calg.Put(id, m)

	return nil
}

// Delete .
func (repo MysqlRepo) Delete(id uint) error {
	if err := repo.db.Delete(nil, "id = ?", id).Error; err != nil {
		return err
	}

	repo.calg.Delete(id)
	return nil
}
```

execute `go run main.go`

```sh
➜  custom-cache-manage git:(master) ✗ go run main.go 
this queryid=9 cost: 245505 ns
this queryid=1 cost: 131838 ns
this queryid=3 cost: 128272 ns
this queryid=2 cost: 112281 ns
this queryid=7 cost: 123942 ns
this queryid=4 cost: 140267 ns
this queryid=7 cost: 148814 ns
this queryid=9 cost: 126904 ns
this queryid=6 cost: 129676 ns
this queryid=2 cost: 174202 ns
this queryid=1 cost: 151673 ns
this queryid=4 cost: 156370 ns
this queryid=3 cost: 159285 ns
this queryid=6 cost: 142215 ns
this queryid=3 cost: 691 ns
this queryid=1 cost: 450 ns
this queryid=8 cost: 160263 ns
this queryid=5 cost: 149655 ns
this queryid=4 cost: 756 ns
this queryid=8 cost: 143363 ns
this queryid=3 cost: 740 ns
this queryid=9 cost: 558 ns
this queryid=2 cost: 476 ns
this queryid=5 cost: 184098 ns
this queryid=1 cost: 824 ns
this queryid=8 cost: 556 ns
this queryid=9 cost: 632 ns
this queryid=7 cost: 480 ns
this queryid=5 cost: 439 ns
this queryid=5 cost: 409 ns
this queryid=7 cost: 431 ns
this queryid=6 cost: 479 ns
this queryid=4 cost: 423 ns
this queryid=8 cost: 423 ns
this queryid=1 cost: 411 ns
this queryid=6 cost: 423 ns
this queryid=8 cost: 394 ns
this queryid=7 cost: 410 ns
this queryid=9 cost: 424 ns
this queryid=4 cost: 428 ns
this queryid=2 cost: 433 ns
this queryid=4 cost: 420 ns
this queryid=9 cost: 424 ns
this queryid=6 cost: 406 ns
this queryid=6 cost: 399 ns
this queryid=5 cost: 405 ns
this queryid=2 cost: 428 ns
this queryid=9 cost: 383 ns
this queryid=4 cost: 399 ns
this queryid=7 cost: 413 ns
this queryid=4 cost: 381 ns
this queryid=1 cost: 427 ns
this queryid=2 cost: 430 ns
this queryid=1 cost: 468 ns
this queryid=1 cost: 406 ns
this queryid=4 cost: 380 ns
this queryid=2 cost: 360 ns
this queryid=3 cost: 660 ns
this queryid=6 cost: 393 ns
this queryid=5 cost: 419 ns
this queryid=7 cost: 1254 ns
this queryid=6 cost: 723 ns
this queryid=4 cost: 503 ns
this queryid=8 cost: 448 ns
this queryid=3 cost: 510 ns
this queryid=1 cost: 432 ns
this queryid=2 cost: 999 ns
this queryid=1 cost: 419 ns
this queryid=8 cost: 658 ns
this queryid=9 cost: 1322 ns
this queryid=9 cost: 543 ns
this queryid=4 cost: 1311 ns
this queryid=5 cost: 348 ns
this queryid=4 cost: 309 ns
this queryid=5 cost: 350 ns
this queryid=9 cost: 311 ns
this queryid=5 cost: 336 ns
this queryid=3 cost: 567 ns
this queryid=9 cost: 293 ns
this queryid=7 cost: 338 ns
this queryid=4 cost: 499 ns
this queryid=7 cost: 318 ns
this queryid=3 cost: 330 ns
this queryid=7 cost: 322 ns
this queryid=6 cost: 339 ns
this queryid=7 cost: 1273 ns
this queryid=4 cost: 1175 ns
this queryid=6 cost: 306 ns
this queryid=2 cost: 316 ns
this queryid=5 cost: 330 ns
this queryid=5 cost: 322 ns
this queryid=6 cost: 324 ns
this queryid=8 cost: 291 ns
this queryid=2 cost: 310 ns
this queryid=3 cost: 321 ns
this queryid=3 cost: 294 ns
this queryid=6 cost: 293 ns
this queryid=8 cost: 3566 ns
```