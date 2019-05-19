package transport

import (
	"github.com/efrengarcial/framework/users/pkg/service"
	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// Server holds the dependencies for a HTTP server.
type Server struct {
	UserService service.UsersService

	Logger kitlog.Logger

	router chi.Router
}

// New returns a new HTTP server.
func New(us service.UsersService,logger kitlog.Logger) *Server {
	s := &Server{
		UserService:  us,
		Logger:   logger,
	}

	r := chi.NewRouter()

	r.Use(accessControl)

	r.Route("/api", func(r chi.Router) {
		h := userHandler{s.UserService, s.Logger}
		r.Mount("/v1", h.router())
	})

	r.Method("GET", "/metrics", promhttp.Handler())

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