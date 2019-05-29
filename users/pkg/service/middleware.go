package service

import (
	"context"
	"github.com/efrengarcial/framework/users/pkg/model"
)

// Middleware describes a service middleware.
type Middleware func(UsersService) UsersService

type authMiddleware struct {
	next UsersService
}

// AuthMiddleware returns a UsersService Middleware.
func AuthMiddleware() Middleware {
	return func(next UsersService) UsersService {
		return &authMiddleware{next}
	}

}
func (a authMiddleware) Create(ctx context.Context, req *model.User) (m0 *model.User, e1 error) {
	// Implement your middleware logic here

	return a.next.Create(ctx, req)
}
