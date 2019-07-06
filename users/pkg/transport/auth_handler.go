package transport

import (
	"context"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/efrengarcial/framework/users/pkg/service"
	"github.com/gin-gonic/gin"
	kitlog "github.com/go-kit/kit/log"
	"net/http"
)


type authHandler struct {
	service service.AuthService
	logger  kitlog.Logger
}


func (h *authHandler) signIn(c *gin.Context) {
	ctx := context.Background()

	loginVM := new(service.LoginVM)
	c.BindJSON(&loginVM)

	token:= new(model.Token)
	err := h.service.Auth(ctx, loginVM, token)

	if err != nil {
		encodeError1(err, h.logger, c)
		return
	}

	c.JSON(http.StatusCreated, token)

}


