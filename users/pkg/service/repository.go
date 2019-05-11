package service

import (
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/efrengarcial/framework/users/pkg/model"
)

// mockery -name=Repository
type Repository interface {

	// Insert puts a new instance of the give Model in the database
	Insert(model model.IModel) (model.IModel, error)

	Update(model model.IModel) error

	Save(model model.IModel) (uint64, error)

	Find(receiver model.IModel, id uint64) error

	FindFirst(receiver model.IModel, where string, args ...interface{}) error

	FindAll(result interface{}, where string, args ...interface{}) (err error)

	FindAllPageable(pageable model.Pageable, result interface{},  where string, args ...interface{}) *pagination.Paginator

	Delete(model model.IModel, where string, args ...interface{}) error

	// NewRecord check if the model exist in the store
	NewRecord(model model.IModel) bool
}