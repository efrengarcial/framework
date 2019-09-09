// Package user defines a User type that's stored in Contexts.
package web

import (
	"context"

	"github.com/efrengarcial/framework/internal/users"
)

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// userKey is the key for user.User values in Contexts. It is
// unexported; clients use user.NewContext and user.FromContext
// instead of using this key directly.
var userKey key

// NewContext returns a new Context that carries value u.
func NewContext(ctx context.Context, u *users.User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext returns the User value stored in ctx, if any.
func FromContext(ctx context.Context) (*users.User, bool) {
	u, ok := ctx.Value(userKey).(*users.User)
	return u, ok
}
