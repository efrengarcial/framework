package service

import (
	"context"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/efrengarcial/framework/users/pkg/utils/paginations"
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


func (a authMiddleware) Update(ctx context.Context, user *model.User) (*model.User, error) {
	panic("implement me")
}


func (a authMiddleware) FindAll(pageable model.Pageable, result interface{}, where string, args ...interface{}) (*paginations.Pagination, error) {
	panic("implement me")
}
