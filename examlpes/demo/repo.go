package main

import (
	"github.com/jinzhu/gorm"
	cp "github.com/yeqown/cached-repository"
)

type userModel struct {
	gorm.Model
	Name     string `gorm:"column:name"`
	Province string `gorm:"column:province"`
	City     string `gorm:"column:city"`
}

// MysqlRepo .
type MysqlRepo struct {
	db *gorm.DB
	*cp.EmbedRepo
}

// GetByID .
func (repo MysqlRepo) GetByID(id uint) (*userModel, error) {
	return nil, nil
}

// Update .
func (repo MysqlRepo) Update(id uint, m *userModel) error {
	return nil
}

// Delete .
func (repo MysqlRepo) Delete(id uint) error {
	return nil
}
