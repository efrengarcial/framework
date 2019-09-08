package repository

import (
	"context"

	base "github.com/efrengarcial/framework/internal/platform/model"
)

// mockery -name=Repository
type Repository interface {

	// Insert puts a new instance of the give Model in the database
	Insert(ctx context.Context, model base.IModel) (base.IModel, error)

	Update(model base.IModel) error

	Save(model base.IModel) (uint64, error)

	Find(receiver base.IModel, id uint64) error

	FindFirst(receiver base.IModel, where string, args ...interface{}) error

	FindAll(result interface{}, where string, args ...interface{}) (err error)

	FindAllPageable(ctx context.Context, pageable *base.Pageable, result interface{},  where string, args ...interface{}) (*base.Pagination, error)

	Delete(model base.IModel) error

	// NewRecord check if the model exist in the store
	NewRecord(model base.IModel) bool
}
