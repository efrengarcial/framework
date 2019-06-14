package service

import (
	"context"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/go-kit/kit/log"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)


// UserService describes the service.
type UserService interface {
	Create(ctx context.Context, req *model.User) (*model.User, error)
}


// service implements the User Service
type userService struct {
	repository UserRepository
	logger     log.Logger
}


// NewService creates and returns a new User service instance
func NewService(rep UserRepository, logger log.Logger) UserService {
	return &userService {
		repository: rep,
		logger:     logger,
	}
}

//https://github.com/pkg/errors
//https://medium.com/@hussachai/error-handling-in-go-a-quick-opinionated-guide-9199dd7c7f76
func (service *userService) Create(ctx context.Context, user *model.User) (*model.User, error) {
	var  (
		err error
		foundUser *model.User
	)

	if user.ID > 0 {
		return nil, NewErrBadRequest("Un nuevo usuario ya no puede tener una ID","userManagement.idexists")
	}

	if  foundUser  , err =  service.repository.FindOneByLogin(strings.ToLower(user.Login)); err == nil && foundUser != nil {
		return nil, NewErrLoginAlreadyUsed("Nombre de inicio de sesión ya usado!", "userManagement.userexists")
	}
	if err != nil { return nil, err}

	if len(user.LangKey) ==0 {
		user.LangKey = "en"
	}

	// Generates a hashed version of our password
	randomPassword, _ := GeneratePassword()
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(randomPassword), bcrypt.DefaultCost)
	if err != nil { return nil, err}

	user.Password = string(hashedPass)
	resetKey, _ := GenerateResetKey()
	user.ResetKey = resetKey
	user.ResetDate = time.Now()
	user.Activated = true

	newUser, err  := service.repository.Insert(user)
	if err != nil { return nil, err}


	return newUser.(*model.User), nil
}


