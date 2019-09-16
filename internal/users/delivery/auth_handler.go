package delivery

import (
	"net/http"

	"github.com/efrengarcial/framework/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


type authHandler struct {
	service users.AuthService
	logger  *logrus.Logger
}


func (h *authHandler) signIn(c *gin.Context) {
	var loginVM users.LoginVM
	if err := c.ShouldBindJSON(&loginVM); err != nil {
		c.Error(err)
		return
	}

	token:= new(users.Token)
	err := h.service.Auth(c.Request.Context(), &loginVM, token)
	if err != nil {
		switch err {
		case users.ErrAuthenticationFailure:
			c.JSON(http.StatusUnauthorized, token)
			return
		default:
			c.Error(err)
			return
		}
	}

	c.JSON(http.StatusOK, token)
}


