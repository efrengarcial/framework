package repository

import (
	"context"
	"github.com/efrengarcial/framework/internal/platform/database"
	"github.com/efrengarcial/framework/internal/platform/service"
)

// mockery -name=Repository
type Repository interface {

	// Insert puts a new instance of the give Model in the database
	Insert(ctx context.Context, model service.IModel) (service.IModel, error)

	Update(model service.IModel) error

	Save(model service.IModel) (uint64, error)

	Find(receiver service.IModel, id uint64) error

	FindFirst(receiver service.IModel, where string, args ...interface{}) error

	FindAll(result interface{}, where string, args ...interface{}) (err error)

	FindAllPageable(ctx context.Context, pageable *service.Pageable, result interface{},  where string, args ...interface{}) (*database.Pagination, error)

	Delete(model service.IModel) error

	// NewRecord check if the model exist in the store
	NewRecord(model service.IModel) bool
}
