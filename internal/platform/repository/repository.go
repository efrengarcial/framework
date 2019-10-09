package repository

import (
	"context"
	"github.com/efrengarcial/framework/internal/domain"
)

// mockery -name=Repository
type Repository interface {

	// Insert puts a new instance of the give Model in the database
	Insert(ctx context.Context, model domain.IModel) error

	Update(model domain.IModel) error

	Save(model domain.IModel) (uint64, error)

	Find(receiver domain.IModel, id uint64) error

	FindFirst(receiver domain.IModel, where string, args ...interface{}) error

	FindAll(result interface{}, where string, args ...interface{}) (err error)

	FindAllPageable(ctx context.Context, pageable *domain.Pageable, result interface{},  where string, args ...interface{}) (*domain.Pagination, error)

	Delete(model domain.IModel) error

	// NewRecord check if the model exist in the store
	NewRecord(model domain.IModel) bool
}
