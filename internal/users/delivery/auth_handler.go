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
	//err :=  c.BindJSON(&loginVM)
	if err := c.ShouldBindJSON(&loginVM); err != nil {
		//c.AbortWithError(400, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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


