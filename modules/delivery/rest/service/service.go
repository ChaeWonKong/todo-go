package service

import (
	"todo-clone/modules/domains"
	"todo-clone/modules/repository"
)

type Service interface {
	FindOne(id uint64) (*domains.Item, error)
	FindAll(offset, limit int) ([]domains.Item, error)
	Insert(item *domains.Item) (int64, error)
	Update(id uint64, title string) (*domains.Item, error)
	Delete(id uint64) (int64, error)
}

type ServiceImpl struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) Service {
	return &ServiceImpl{repo}
}

func (s ServiceImpl) FindOne(id uint64) (item *domains.Item, err error) {
	item = &domains.Item{}
	tx := s.repo.First(item, id)
	if tx.Error != nil || tx.RowsAffected == 0 {
		item = nil
	}

	return item, tx.Error
}

func (s ServiceImpl) FindAll(offset, limit int) ([]domains.Item, error) {
	items := make([]domains.Item, limit)
	err := s.repo.Limit(limit).Offset(offset).Find(&items).Error

	return items, err
}

func (s ServiceImpl) Insert(item *domains.Item) (int64, error) {
	tx := s.repo.Create(item)
	return tx.RowsAffected, tx.Error
}

func (s ServiceImpl) Update(id uint64, title string) (item *domains.Item, err error) {
	item = &domains.Item{}
	tx := s.repo.Model(item).Where("id = ?", id).Update("title", title)

	if tx.Error != nil || tx.RowsAffected == 0 {
		item = nil
	}

	return item, tx.Error
}

func (s ServiceImpl) Delete(id uint64) (int64, error) {
	tx := s.repo.Delete(&domains.Item{ID: id})
	return tx.RowsAffected, tx.Error
}
