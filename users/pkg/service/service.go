package service

import (
	"context"
	"github.com/efrengarcial/framework/users/pkg/model"
)

// UsersService describes the service.
type UsersService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
	Create(ctx context.Context, req *model.User) (*model.User, error)
}
