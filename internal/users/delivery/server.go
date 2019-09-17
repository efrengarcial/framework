package delivery

import (
	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/efrengarcial/framework/internal/mid"
	"github.com/efrengarcial/framework/internal/platform/auth"
	"github.com/efrengarcial/framework/internal/users"
	"go.opencensus.io/plugin/ochttp"
	"net/http"
	"os"

	"github.com/efrengarcial/framework/internal/platform/web"
	"github.com/efrengarcial/framework/internal/users/repository"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
)

func setUserRouter(api *gin.RouterGroup, us users.UserService, logger *logrus.Logger, authenticator *auth.Authenticator) {

	api.Use(mid.Authenticate(authenticator))
	{
		api.Use(mid.HasRole(auth.RoleAdmin))
		h := userHandler{us, logger}
		api.POST("/users", func(c *gin.Context) {
			ochttp.SetRoute(c.Request.Context(), "/users")
		}, h.createUser)
		api.PUT("/users", func(c *gin.Context) {
			ochttp.SetRoute(c.Request.Context(), "/users")
		}, h.updateUser)
		api.GET("/users", func(c *gin.Context) {
			ochttp.SetRoute(c.Request.Context(), "/users")
		}, h.findAll)
	}
}

func setAuthRouter(router *gin.Engine, as users.AuthService, logger *logrus.Logger) {
	a := authHandler{as, logger}
	router.POST("/api/authenticate", a.signIn)
}

//New returns a new HTTP server.
func New(shutdown chan os.Signal, db *gorm.DB, logger *logrus.Logger, exporter *prometheus.Exporter, authenticator *auth.Authenticator) http.Handler  {
	// Setup repositories
	repo := repository.NewUserGormRepository(db)
	us := users.NewService(repo, logger)
	as := users.NewAuthService(repo, authenticator, logger)

	app := web.NewApp(shutdown, logger)
	router := app.Engine
	router.Use(mid.Error(logger))
	v1 := router.Group("/api/v1")
	setUserRouter(v1, us, logger, authenticator)
	setAuthRouter(router, as, logger)

	router.GET("/metrics", func(c *gin.Context) {
		exporter.ServeHTTP(c.Writer, c.Request)
	})

	return app
}
