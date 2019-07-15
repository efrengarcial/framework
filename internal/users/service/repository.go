package service

import (
	"github.com/efrengarcial/framework/internal/platform/database"
)

// mockery -name=Repository
type Repository interface {

	// Insert puts a new instance of the give Model in the database
	Insert(model IModel) (IModel, error)

	Update(model IModel) error

	Save(model IModel) (uint64, error)

	Find(receiver IModel, id uint64) error

	FindFirst(receiver IModel, where string, args ...interface{}) error

	FindAll(result interface{}, where string, args ...interface{}) (err error)

	FindAllPageable(pageable *Pageable, result interface{},  where string, args ...interface{}) (*database.Pagination, error)

	Delete(model IModel, where string, args ...interface{}) error

	// NewRecord check if the model exist in the store
	NewRecord(model IModel) bool
}
