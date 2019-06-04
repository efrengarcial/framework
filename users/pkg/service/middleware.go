package service

import (
	"context"
	"github.com/efrengarcial/framework/users/pkg/model"
)

// Middleware describes a service middleware.
type Middleware func(UserService) UserService

type authMiddleware struct {
	next UserService
}

// AuthMiddleware returns a UserService Middleware.
func AuthMiddleware() Middleware {
	return func(next UserService) UserService {
		return &authMiddleware{next}
	}

}
func (a authMiddleware) Create(ctx context.Context, req *model.User) (m0 *model.User, e1 error) {
	// Implement your middleware logic here

	return a.next.Create(ctx, req)
}
