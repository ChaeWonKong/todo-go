package repository

// import (
// 	"todo-clone/modules/domains"

// 	"gorm.io/gorm"
// )

// type Finder interface {
// 	Find(offset, limit int) *gorm.DB
// }

// func NewFinder(repo *Repository) Finder {
// 	return repo
// }

// func (repo Repository) Find(offset, limit int) *gorm.DB {
// 	items := make([]domains.Item, limit)
// 	tx := repo.DB.Limit(limit).Offset(offset).Find(&items)

// 	return tx
// }
