package transport

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/efrengarcial/framework/users/pkg/service"
	"github.com/go-chi/chi"
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

	router chi.Router
}

// New returns a new HTTP server.
func New(us UserService,as AuthService, logger kitlog.Logger) *Server {
	s := &Server{
		UserService:  us,
		AuthService: as,
		logger:   logger,
	}

	r := chi.NewRouter()

	r.Use(accessControl)

	r.Route("/api", func(r chi.Router) {
		h := userHandler{s.UserService, s.logger}
		r.Mount("/v1", h.router())

		a := authHandler{s.AuthService, s.logger}
		r.Mount("/authenticate", a.router())
	})

	s.router = r

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

type iErrBadRequest interface {
	GetErrorKey() string
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func encodeError(_ context.Context, err error,  logger kitlog.Logger, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var status int

	switch err.(type) {
	case iErrBadRequest:
		w.WriteHeader(http.StatusBadRequest)
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
		w.WriteHeader(http.StatusInternalServerError)
		status = http.StatusInternalServerError
	}


	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": err.Error(),
		"status" : status,
	})
}
