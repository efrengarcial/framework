package delivery

import (
	"github.com/efrengarcial/framework/internal/domain"
	"net/http"

	"github.com/efrengarcial/framework/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


type authHandler struct {
	service user.AuthService
	logger  *logrus.Logger
}


func (h *authHandler) signIn(c *gin.Context) {
	var loginVM domain.LoginVM
	if err := c.ShouldBindJSON(&loginVM); err != nil {
		c.Error(err)
		return
	}

	token:= new(domain.Token)
	err := h.service.Auth(c.Request.Context(), &loginVM, token)
	if err != nil {
		switch err {
		case domain.ErrAuthenticationFailure:
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "error.http.401", "title" : "Unauthorized",  "detail" : "Bad credentials" ,
					"path" : c.Request.URL.Path, "status" : http.StatusUnauthorized})
			return
		default:
			c.Error(err)
			return
		}
	}

	c.JSON(http.StatusOK, token)
}


