package handlers

import (
	"fmt"
	"github.com/efrengarcial/framework/internal/users/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

// Server holds the dependencies for a HTTP server.
type Server struct {
	UserService service.UserService
	AuthService service.AuthService

	logger kitlog.Logger

	router *gin.Engine
}

func SetupUserRouter(us service.UserService, logger kitlog.Logger) *gin.Engine {

	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.Default())

	v1 := router.Group("/api/v1")
	{
		h := userHandler{us, logger}
		v1.POST("/users" , h.createUser)
		v1.PUT("/users", h.updateUser)
		v1.GET("/users", h.findAll)
	}

	return router
}

func setupRouter(us service.UserService,as service.AuthService, logger kitlog.Logger) *gin.Engine {

	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.Default())

	v1 := router.Group("/api/v1")
	{
		h := userHandler{us, logger}
		v1.POST("/users" , h.createUser)
		v1.PUT("/users", h.updateUser)
		v1.GET("/users", h.findAll)
	}

	a := authHandler{as, logger}
	router.POST("/authenticate", a.signIn)
	return router
}

// New returns a new HTTP server.
func New(us service.UserService,as service.AuthService, logger kitlog.Logger) http.Handler  {
	s := &Server{
		UserService:  us,
		AuthService: as,
		logger:   logger,
	}

	s.router = setupRouter(us, as, logger)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

type iErrBadRequest interface {
	GetErrorKey() string
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func encodeError(err error,  logger kitlog.Logger,c *gin.Context) {
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
		level.Error(logger).Log("error", errorLog.String())
		fmt.Printf("with stack trace => %+v \n\n", err)
		status = http.StatusInternalServerError
	}

	c.AbortWithStatusJSON(status, gin.H{
		"message": err.Error(),
		"status" : status,
	} )
}
