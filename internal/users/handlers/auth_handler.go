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
	var loginVM service.LoginVM
	//err :=  c.BindJSON(&loginVM)
	if err := c.ShouldBindJSON(&loginVM); err != nil {
		//c.AbortWithError(400, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token:= new(service.Token)
	err := h.service.Auth(c.Request.Context(), &loginVM, token)
	if err != nil {
		switch err {
		case service.ErrAuthenticationFailure:
			c.JSON(http.StatusUnauthorized, token)
			return
		default:
			c.Error(err)
			return
		}
	}

	c.JSON(http.StatusOK, token)

}


