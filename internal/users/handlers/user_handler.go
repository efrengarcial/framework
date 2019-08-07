package handlers

import (
	base "github.com/efrengarcial/framework/internal/platform/service"
	"github.com/efrengarcial/framework/internal/users/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type userHandler struct {
	service service.UserService
	logger  *logrus.Logger
}

func (h *userHandler) createUser(c *gin.Context) {

	var user *service.User
	c.BindJSON(&user)

	user, err := h.service.Create(c.Request.Context(),user)
	if err != nil {
		encodeError( err, h.logger, c)
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

	var user *service.User
	c.BindJSON(&user)

	user, err := h.service.Update(c.Request.Context(), user)
	if err != nil {
		encodeError( err, h.logger, c)
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
	pageable :=  new(base.Pageable)
	page, err := strconv.Atoi(c.Query("page"))
	limit, err := strconv.Atoi(c.Query("limit"))
	pageable.Page = page
	pageable.Limit = limit
	pageable.OrderBy= c.QueryArray("orderBy")

	var users []service.User
	//pageable := model.Pageable{Page:1 , Limit: 10 , OrderBy: []string{"id desc"} , ShowSQL:true}
	_, err = h.service.FindAll(c.Request.Context(), pageable, &users, "")
	if err != nil {
		encodeError( err, h.logger, c)
		return
	}

	var response = struct { Users []service.User `json:"users"` }{
		Users: users,
	}

	c.JSON(http.StatusOK, response)
}
