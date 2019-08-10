package handlers

import (
	"github.com/efrengarcial/framework/internal/users/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)


type authHandler struct {
	service service.AuthService
	logger  *logrus.Logger
}


func (h *authHandler) signIn(c *gin.Context) {
	loginVM := new(service.LoginVM)
	c.BindJSON(&loginVM)

	token:= new(service.Token)
	err := h.service.Auth(c.Request.Context(), loginVM, token)

	if err != nil {
		encodeError(err, h.logger, c)
		return
	}

	c.JSON(http.StatusCreated, token)

}


