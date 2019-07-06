package transport

import (
	"fmt"
	. "github.com/efrengarcial/framework/users/pkg/service"
	"github.com/gin-gonic/gin"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

// Server holds the dependencies for a HTTP server.
type Server struct {
	UserService UserService
	AuthService AuthService

	logger kitlog.Logger

	router *gin.Engine
}

// New returns a new HTTP server.
func New(us UserService,as AuthService, logger kitlog.Logger) *Server {
	s := &Server{
		UserService:  us,
		AuthService: as,
		logger:   logger,
	}

	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(corsMiddleware())

	v1 := router.Group("/api/v1")
	{
		h := userHandler{s.UserService, s.logger}
		v1.POST("/users" , h.createUser)
		v1.PUT("/users", h.updateUser)
		v1.GET("/users", h.findAll)
	}

	a := authHandler{s.AuthService, s.logger}
	router.POST("/authenticate", a.signIn)

	s.router = router

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

type iErrBadRequest interface {
	GetErrorKey() string
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func encodeError1(err error,  logger kitlog.Logger,c *gin.Context) {
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
