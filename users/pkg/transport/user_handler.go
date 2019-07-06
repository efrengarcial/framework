package transport

import (
	"context"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/efrengarcial/framework/users/pkg/service"
	"github.com/gin-gonic/gin"
	kitlog "github.com/go-kit/kit/log"
	"net/http"
)

type userHandler struct {
	service service.UserService
	logger  kitlog.Logger
}

func (h *userHandler) createUser(c *gin.Context) {

	ctx := context.Background()
	var user *model.User
	c.BindJSON(&user)

	user, err := h.service.Create(ctx,user)
	if err != nil {
		encodeError1( err, h.logger, c)
		return
	}

	var response = struct {
		ID uint64 `json:"id"`
	}{
		ID: user.GetID(),
	}

	c.JSON(http.StatusCreated, response)
}


func (h *userHandler) updateUser(c *gin.Context) {

	ctx := context.Background()
	var user *model.User
	c.BindJSON(&user)

	user, err := h.service.Update(ctx, user)
	if err != nil {
		encodeError1( err, h.logger, c)
		return
	}

	var response = struct {
		ID uint64 `json:"id"`
	}{
		ID: user.GetID(),
	}

	c.JSON(http.StatusOK, response)
}


func (h *userHandler) findAll(c *gin.Context) {
	var users []model.User
	pageable := model.Pageable{Model: &model.User{}, Page:1 , Limit: 10 , OrderBy: []string{"id desc"} , ShowSQL:true}
	_, err := h.service.FindAll(pageable, &users, "")
	if err != nil {
		encodeError1( err, h.logger, c)
		return
	}

	var response = struct { Users []model.User `json:"users"` }{
		Users: users,
	}

	c.JSON(http.StatusOK, response)
}
