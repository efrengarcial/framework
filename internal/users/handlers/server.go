package handlers

import (
	"fmt"
	"go.opencensus.io/plugin/ochttp"
	"net/http"
	"os"
	"strings"

	"github.com/efrengarcial/framework/internal/platform/web"
	"github.com/efrengarcial/framework/internal/users/repository"
	"github.com/efrengarcial/framework/internal/users/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func setUserRouter(router *gin.Engine, us service.UserService, logger *logrus.Logger) {

	v1 := router.Group("/api/v1")
	{
		h := userHandler{us, logger}
		v1.POST("/users" , func(c *gin.Context) {
			ochttp.SetRoute(c.Request.Context(), "/users")
		}, h.createUser)
		v1.PUT("/users",func(c *gin.Context) {
			ochttp.SetRoute(c.Request.Context(), "/users")
		}, h.updateUser)
		v1.GET("/users",func(c *gin.Context) {
			ochttp.SetRoute(c.Request.Context(), "/users")
		}, h.findAll)
	}
}

func setAuthRouter(router *gin.Engine, as service.AuthService, logger *logrus.Logger) {
	a := authHandler{as, logger}
	router.POST("/authenticate", a.signIn)
}

//New returns a new HTTP server.
func New(shutdown chan os.Signal, db *gorm.DB, logger *logrus.Logger) http.Handler  {
	// Setup repositories
	repo := repository.NewUserGormRepository(db)
	us := service.NewService(repo, logger)
	ts := service.NewTokenService()
	as := service.NewAuthService(repo, ts, logger)

	app := web.NewApp(shutdown, logger)
	router := app.Engine
	setUserRouter(router, us, logger)
	setAuthRouter(router, as, logger)

	return app
}

type iErrBadRequest interface {
	GetErrorKey() string
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func encodeError(err error,  logger *logrus.Logger,c *gin.Context) {
	var status int

	switch err.(type) {
	case iErrBadRequest:
		status = http.StatusBadRequest
		iError, _ := err.(iErrBadRequest)
		fmt.Println(iError.GetErrorKey())
	default:
		var errorLog strings.Builder
		errorLog.WriteString(err.Error() + "\n\n")
		if err, ok := err.(stackTracer); ok {
			for _, f := range err.StackTrace() {
				errorLog.WriteString(fmt.Sprintf("%+v \n\n", f))
			}
		}
		logger.Error("error", errorLog.String())
		fmt.Printf("with stack trace => %+v \n\n", err)
		status = http.StatusInternalServerError
	}

	c.AbortWithStatusJSON(status, gin.H{
		"message": err.Error(),
		"status" : status,
	} )
}
