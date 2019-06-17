package service

import "github.com/efrengarcial/framework/users/pkg/model"

type UserRepository interface {
	Repository
	GetByEmail(email string) (*model.User, error)
	GetByLogin(login string) (*model.User, error)
	FindOneByLogin(login string) (*model.User, error)
	FindOneByEmail(login string) (*model.User, error)
}
