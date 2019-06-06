package transport

import (
	"github.com/efrengarcial/framework/users/pkg/service"
	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
	"net/http"
)

// Server holds the dependencies for a HTTP server.
type Server struct {
	UserService service.UserService
	AuthService service.AuthService

	Logger kitlog.Logger

	router chi.Router
}

// New returns a new HTTP server.
func New(us service.UserService,as service.AuthService, logger kitlog.Logger) *Server {
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
	})

	r.Route("/api", func(r chi.Router) {
		h := authHandler{s.AuthService, s.Logger}
		r.Mount("/", h.router())
	})

	r.Mount("/login", adminRouter())

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
