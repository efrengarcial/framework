package delivery

import (
	"net/http"
	"strconv"

	base "github.com/efrengarcial/framework/internal/platform/model"
	"github.com/efrengarcial/framework/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type userHandler struct {
	service users.UserService
	logger  *logrus.Logger
}

func (h *userHandler) createUser(c *gin.Context) {

	var user *users.User
	//c.BindJSON(&user)
	if err := c.ShouldBindJSON(&user); err != nil {
		//c.AbortWithError(400, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Create(c.Request.Context(), user)
	if err != nil {
		c.Error(err)
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

	var user *users.User
	c.BindJSON(&user)

	user, err := h.service.Update(c.Request.Context(), user)
	if err != nil {
		c.Error(err)
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

	var usersList []users.User
	//pageable := model.Pageable{Page:1 , Limit: 10 , OrderBy: []string{"id desc"} , ShowSQL:true}
	_, err = h.service.FindAll(c.Request.Context(), pageable, &usersList, "")
	if err != nil {
		c.Error(err)
		return
	}

	var response = struct { Users []users.User `json:"users"` }{
		Users: usersList,
	}

	c.JSON(http.StatusOK, response)
}
