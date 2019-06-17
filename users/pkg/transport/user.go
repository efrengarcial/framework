package transport

import (
	"context"
	"encoding/json"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/efrengarcial/framework/users/pkg/service"
	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"net/http"
)

type userHandler struct {
	service service.UserService
	logger  kitlog.Logger
}

func (h *userHandler) router() chi.Router {
	r := chi.NewRouter()

	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.createUser)
	})

	return r
}

func (h *userHandler) createUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var user = new(model.User)

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		level.Error(h.logger).Log("error", err)
		encodeError(ctx, err, h.logger, w)
		return
	}

	user, err := h.service.Create(ctx,user)
	if err != nil {
		encodeError(ctx, err, h.logger, w)
		return
	}

	var response = struct {
		ID uint64 `json:"id"`
	}{
		ID: user.GetID(),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		level.Error(h.logger).Log("error", err)
		encodeError(ctx, err, h.logger, w)
		return
	}
}
