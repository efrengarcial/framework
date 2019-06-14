package transport

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/efrengarcial/framework/users/pkg/service"
	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
	"log"
	"net/http"
)

// Server holds the dependencies for a HTTP server.
type Server struct {
	UserService UserService
	AuthService AuthService

	Logger kitlog.Logger

	router chi.Router
}

// New returns a new HTTP server.
func New(us UserService,as AuthService, logger kitlog.Logger) *Server {
	s := &Server{
		UserService:  us,
		AuthService: as,
		Logger:   logger,
	}

	r := chi.NewRouter()

	r.Use(accessControl)

	r.Route("/api", func(r chi.Router) {
		h := userHandler{s.UserService, s.Logger}
		r.Mount("/v1", h.router())

		a := authHandler{s.AuthService, s.Logger}
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


func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var status int

	switch err.(type) {
	case iErrBadRequest:
		w.WriteHeader(http.StatusBadRequest)
		status = http.StatusBadRequest
		iError, _ := err.(iErrBadRequest)
		fmt.Println(iError.GetErrorKey())
	default:
		w.WriteHeader(http.StatusInternalServerError)
		status = http.StatusInternalServerError
		log.Print(err)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": err.Error(),
		"status" : status,
	})
}
