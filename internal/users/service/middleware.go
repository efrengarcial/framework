package service

import (
	"context"
	"github.com/efrengarcial/framework/internal/platform/database"
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
func (a authMiddleware) Create(ctx context.Context, req *User) (m0 *User, e1 error) {
	// Implement your middleware logic here

	return a.Create(ctx, req)
}


func (a authMiddleware) Update(ctx context.Context, user *User) (*User, error) {
	panic("implement me")
}


func (a authMiddleware) FindAll(pageable *Pageable, result interface{}, where string, args ...interface{}) (*database.Pagination, error) {
	panic("implement me")
}
