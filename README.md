# cached-repository
a basic repository to support cached data in memory and is based LRU-K cache replacing algorithm

### TODOs

* [x] LRU-1 & LRU-K

* [x] Cached-Repository [demo](./examples/custom-cache-manage/main.go)

* [ ] concurrent safe

### Demo

```go
package main

import (
	"fmt"
	"math/rand"
	"time"

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

	for i := 0; i < 10; i++ {
		repo.Create(&userModel{
			Model: gorm.Model{
				ID: uint(i + 1),
			},
			Name:     fmt.Sprintf("name-%d", i+1),
			Province: fmt.Sprintf("province-%d", i+1),
			City:     fmt.Sprintf("city-%d", i+1),
		})
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000; i++ {
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
	}
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
		fmt.Printf("this query cost: %d nano seconds\n", time.Now().Sub(start).Nanoseconds())
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
this query cost: 269226 nano seconds
this query cost: 133925 nano seconds
this query cost: 168136 nano seconds
this query cost: 135615 nano seconds
this query cost: 141056 nano seconds
this query cost: 127871 nano seconds
this query cost: 128350 nano seconds
this query cost: 127124 nano seconds
this query cost: 129247 nano seconds
this query cost: 121573 nano seconds
this query cost: 456 nano seconds
this query cost: 442 nano seconds
this query cost: 117046 nano seconds
this query cost: 117790 nano seconds
this query cost: 485 nano seconds
this query cost: 327 nano seconds
this query cost: 329 nano seconds
this query cost: 323 nano seconds
this query cost: 308 nano seconds
this query cost: 126181 nano seconds
this query cost: 461 nano seconds
this query cost: 399 nano seconds
this query cost: 344 nano seconds
this query cost: 3189 nano seconds
this query cost: 131464 nano seconds
this query cost: 443 nano seconds
this query cost: 340 nano seconds
this query cost: 254 nano seconds
this query cost: 119095 nano seconds
this query cost: 114699 nano seconds
this query cost: 488 nano seconds
this query cost: 346 nano seconds
this query cost: 390 nano seconds
this query cost: 112766 nano seconds
this query cost: 416 nano seconds
this query cost: 316 nano seconds
this query cost: 364 nano seconds
this query cost: 351 nano seconds
this query cost: 291 nano seconds
this query cost: 295 nano seconds
this query cost: 355 nano seconds
this query cost: 321 nano seconds
this query cost: 304 nano seconds
this query cost: 327 nano seconds
this query cost: 391 nano seconds
this query cost: 366 nano seconds
this query cost: 310 nano seconds
this query cost: 323 nano seconds
this query cost: 2603 nano seconds
this query cost: 304464 nano seconds
this query cost: 1002 nano seconds
this query cost: 1174 nano seconds
this query cost: 1130 nano seconds
this query cost: 527 nano seconds
this query cost: 339 nano seconds
this query cost: 638 nano seconds
this query cost: 305 nano seconds
this query cost: 298 nano seconds
... more ignored
```