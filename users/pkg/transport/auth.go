package transport

import (
	"encoding/json"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/efrengarcial/framework/users/pkg/service"
	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
	"net/http"
)


type authHandler struct {
	service service.AuthService
	logger  kitlog.Logger
}

func (h *authHandler) router() http.Handler  {
	r := chi.NewRouter()
	r.Post("/", h.signIn)
	return r
}

func (h *authHandler) signIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	loginVM := new(service.LoginVM)

	if err := json.NewDecoder(r.Body).Decode(&loginVM); err != nil {
		h.logger.Log("error", err)
		encodeError(ctx, err, w)
		return
	}

	token:= new(model.Token)
	err := h.service.Auth(ctx, loginVM, token)

	if err != nil {
		encodeError(ctx, err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(token); err != nil {
		h.logger.Log("error", err)
		encodeError(ctx, err, w)
		return
	}

}

