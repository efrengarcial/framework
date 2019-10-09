package mid

import (
	"github.com/efrengarcial/framework/internal/platform/auth"
	"github.com/efrengarcial/framework/internal/platform/web"
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
			//expected authorization header format: Bearer <token>
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "error.http.401", "title" : "Unauthorized",  "detail" : "Bad credentials" ,
					"path" : c.Request.URL.Path, "status" : http.StatusUnauthorized})
			return
		}

		claims, err := authenticator.ParseClaims(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "error.http.401", "title" : "Unauthorized",  "detail" : "Bad credentials" ,
					"path" : c.Request.URL.Path, "status" : http.StatusUnauthorized})
			return
		}
		// Add claims to the context so they can be retrieved later.
		ctx, err := authenticator.SetUser(c.Request.Context(), claims.Subject)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "error.http.401", "title" : "Unauthorized",  "detail" : "Bad credentials" ,
					"path" : c.Request.URL.Path, "status" : http.StatusUnauthorized})
			return
		}

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// HasPermission validates that an authenticated user has at least one Permission from a
// specified list. This method constructs the actual function that is used.
func HasPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		claims, ok := web.FromContext(c.Request.Context())
		if !ok {
			//claims missing from context: HasPermission called without/before Authenticate
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "error.http.401", "title" : "Unauthorized",  "detail" : "Bad credentials" ,
					"path" : c.Request.URL.Path, "status" : http.StatusUnauthorized})
			return
		}

		if !claims.HasPermission(permissions...) {
			//you are not authorized for that action
			c.AbortWithStatusJSON(http.StatusForbidden,
				gin.H{"error": "error.http.403", "title" : "Forbidden",  "detail" : "Access is denied" ,
					"path" : c.Request.URL.Path, "status" : http.StatusForbidden})
			return
		}
	}
}
