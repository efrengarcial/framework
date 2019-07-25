package service

import "github.com/efrengarcial/framework/internal/platform/repository"

//mockery -name=UserRepository
type UserRepository interface {
	repository.Repository
	GetByEmail(email string) (*User, error)
	GetByLogin(login string) (*User, error)
	FindOneByLogin(login string) (*User, error)
	FindOneByEmail(login string) (*User, error)
}
