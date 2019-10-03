package users

import (
	"context"
	"github.com/efrengarcial/framework/internal/domain"

	"github.com/efrengarcial/framework/internal/platform/repository"
)

//mockery -name=UserRepository
type UserRepository interface {
	repository.Repository
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
	FindOneByLogin(ctx context.Context, login string) (*domain.User, error)
	FindOneByEmail(ctx context.Context, login string) (*domain.User, error)
}
