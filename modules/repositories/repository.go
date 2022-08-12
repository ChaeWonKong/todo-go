package repositories

import (
	"log"
	"todo-go/modules/domains"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	// Create(value interface{}) (tx *gorm.DB)
	// First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	// Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
	// Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	AutoMigrate(dst ...interface{}) error
	// Model(value interface{}) (tx *gorm.DB)

	FindOne(id uint64) (interface{}, error)
	FindAll() (interface{}, error)
	CreateOne(item *domains.Item) (int64, error)
	UpdateOne(id uint64, updateDto map[string]interface{}) (domains.Item, error)
	DeleteOne(id uint64) (int64, error)
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

func (repo RepositoryImpl) FindOne(id uint64) (interface{}, error) {
	item := domains.Item{}

	if tx := repo.First(&item, id); tx.Error != nil {
		return nil, tx.Error
	}

	return item, nil
}

func (repo RepositoryImpl) FindAll() (interface{}, error) {
	items := make([]domains.Item, 1)

	if err := repo.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (repo RepositoryImpl) CreateOne(item *domains.Item) (int64, error) {
	tx := repo.Create(item)

	return tx.RowsAffected, tx.Error
}

func (repo RepositoryImpl) UpdateOne(id uint64, updateDto map[string]interface{}) (domains.Item, error) {
	item := domains.Item{}
	tx := repo.Model(&item).Clauses(clause.Returning{}).Where("id = ?", id).Updates(updateDto)

	return item, tx.Error
}

func (repo RepositoryImpl) DeleteOne(id uint64) (int64, error) {
	item := domains.Item{}
	tx := repo.Delete(&item, id)

	return tx.RowsAffected, tx.Error
}
