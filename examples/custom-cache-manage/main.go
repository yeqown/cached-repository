package main

import (
	"fmt"
	"math/rand"
	"sync"
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
		repo.Create(&UserModel{
			Model: gorm.Model{
				ID: uint(i + 1),
			},
			Name:     fmt.Sprintf("name-%d", i+1),
			Province: fmt.Sprintf("province-%d", i+1),
			City:     fmt.Sprintf("city-%d", i+1),
		})
	}

	wg := sync.WaitGroup{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			id := uint(rand.Intn(10))
			if id == 0 {
				return
			}

			v, err := repo.GetByID(id)
			if err != nil {
				fmt.Printf("err: %d , %v\n", id, err)
				return
			}

			if v.ID != id ||
				v.Name != fmt.Sprintf("name-%d", id) ||
				v.Province != fmt.Sprintf("province-%d", id) ||
				v.City != fmt.Sprintf("city-%d", id) {
				fmt.Printf("err: not matched target with id[%d]: %v\n", v.ID, v)
			}
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

	db.DropTableIfExists(&UserModel{})
	db.AutoMigrate(&UserModel{})

	return NewMysqlRepo(db)
}

// UserModel struct .
type UserModel struct {
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
func (repo MysqlRepo) Create(m *UserModel) error {
	return repo.db.Create(m).Error
}

// GetByID .
func (repo MysqlRepo) GetByID(id uint) (*UserModel, error) {
	start := time.Now()
	defer func() {
		fmt.Printf("this queryid=%d cost: %d ns\n", id, time.Now().Sub(start).Nanoseconds())
	}()

	v, ok := repo.calg.Get(id)
	if ok {
		return v.(*UserModel), nil
	}

	// actual find in DB
	m := new(UserModel)
	if err := repo.db.Where("id = ?", id).First(m).Error; err != nil {
		return nil, err
	}

	repo.calg.Put(id, m)
	return m, nil
}

// Update .
func (repo MysqlRepo) Update(id uint, m *UserModel) error {
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
