package service

import (
	"context"
	"github.com/efrengarcial/framework/internal/platform/repository"
)

//mockery -name=UserRepository
type UserRepository interface {
	repository.Repository
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByLogin(ctx context.Context, login string) (*User, error)
	FindOneByLogin(ctx context.Context, login string) (*User, error)
	FindOneByEmail(ctx context.Context, login string) (*User, error)
}
