package service

import "github.com/efrengarcial/framework/users/pkg/model"

// mockery -name=Stringer
type Repository interface {

	// Insert puts a new instance of the give Model in the database
	Insert(model model.IModel) (model.IModel, error)

	Update(model model.IModel) (error)

	Save(model model.IModel) (uint64, error)

	FindById(receiver model.IModel, uint uint64) (error)

	FindFirst(receiver model.IModel, where string, args ...interface{}) (error)

	FindAll(models interface{}, where string, args ...interface{}) (err error)

	Delete(model model.IModel, where string, args ...interface{}) error

	// NewRecord check if the model exist in the store
	NewRecord(model model.IModel) bool
}