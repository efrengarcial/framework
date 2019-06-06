package transport

import (
	"encoding/json"
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
	r.Post("/authenticate", h.signIn)
	return r
}

func (h *authHandler) signIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data = new(struct {
		UserName    string `json:"username"`
		Password 	string `json:"password"`
		RememberMe 	bool `json:"rememberMe"`
	})

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.logger.Log("error", err)
		encodeError(ctx, err, w)
		return
	}

	h.service.Auth(ctx, )
}

