package mid

import (
	"context"
	"errors"
	"github.com/efrengarcial/framework/internal/platform/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Authenticate(authenticator *auth.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the authorization header. Expected header is of
		// the format `Bearer <token>`.
		parts := strings.Split(c.Request.Header.Get("Authorization"), " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			err := errors.New("expected authorization header format: Bearer <token>")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, err := authenticator.ParseClaims(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Add claims to the context so they can be retrieved later.
		ctx := context.WithValue(c.Request.Context(), auth.Key, claims)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// HasRole validates that an authenticated user has at least one Role from a
// specified list. This method constructs the actual function that is used.
func HasRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		claims, ok := c.Request.Context().Value(auth.Key).(auth.Claims)
		if !ok {
			err :=  errors.New("claims missing from context: HasRole called without/before Authenticate")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if !claims.HasRole(roles...) {
			//you are not authorized for that action
			c.AbortWithStatusJSON(http.StatusForbidden,
				gin.H{"error": "error.http.403", "title" : "Forbidden",  "detail" : "Access is denied" ,
					"path" : c.Request.URL.Path, "status" : http.StatusForbidden})
			return
		}
	}
}
