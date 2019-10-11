package user

import (
	"context"
	"github.com/efrengarcial/framework/internal/domain"
)

// Service describes the User service.
type Service interface {
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User)  error
	FindAll(ctx context.Context, pageable *domain.Pageable, result interface{},  where string, args ...interface{})(*domain.Pagination, error)
}
