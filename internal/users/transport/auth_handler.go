package transport

import (
	"context"
	service2 "github.com/efrengarcial/framework/internal/users/service"
	"github.com/gin-gonic/gin"
	kitlog "github.com/go-kit/kit/log"
	"net/http"
)


type authHandler struct {
	service service2.AuthService
	logger  kitlog.Logger
}


func (h *authHandler) signIn(c *gin.Context) {
	ctx := context.Background()

	loginVM := new(service2.LoginVM)
	c.BindJSON(&loginVM)

	token:= new(service2.Token)
	err := h.service.Auth(ctx, loginVM, token)

	if err != nil {
		encodeError1(err, h.logger, c)
		return
	}

	c.JSON(http.StatusCreated, token)

}


