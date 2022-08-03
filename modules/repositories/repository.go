package repositories

import (
	"log"

	"gorm.io/gorm"
)

type Repository interface {
	Create(value interface{}) (tx *gorm.DB)
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
	Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	AutoMigrate(dst ...interface{}) error
}

type RepositoryImpl struct {
	*gorm.DB
}

func NewRepository(sqliteDialector SQLiteDialector) Repository {
	var dialector gorm.Dialector = sqliteDialector

	db, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	return RepositoryImpl{db}
}
