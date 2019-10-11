package user

import (
	"context"
	"github.com/efrengarcial/framework/internal/domain"
)

// AuthService has the logic authentication
type AuthService interface {
	Auth(ctx context.Context, req *domain.LoginVM, res *domain.Token) error
}
