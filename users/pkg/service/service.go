package service

import (
	"context"
	"errors"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrHashingPassword = errors.New("error hashing password")
	ErrCmdRepository   = errors.New("unable to command repository")
	ErrQueryRepository = errors.New("unable to query repository")
)


// UsersService describes the service.
type UsersService interface {
	Create(ctx context.Context, req *model.User) (*model.User, error)
}


// service implements the Order Service
type userService struct {
	repository Repository
	logger     log.Logger
}


// NewService creates and returns a new Order service instance
func NewService(rep Repository, logger log.Logger) UsersService {
	return &userService {
		repository: rep,
		logger:     logger,
	}
}


func (service *userService) Create(ctx context.Context, user *model.User) (*model.User, error) {
	logger := log.With(service.logger, "method", "Create")
	if len(user.LangKey) ==0 {
		user.LangKey = "en"
	}

	// Generates a hashed version of our password
	randomPassword, _ := GeneratePassword()
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(randomPassword), bcrypt.DefaultCost)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, ErrHashingPassword
	}
	user.Password = string(hashedPass)
	resetKey, _ := GenerateResetKey()
	user.ResetKey = resetKey
	user.ResetDate = time.Now()
	user.Activated = true

	newUser, err  := service.repository.Insert(user)
	if  err != nil {
		level.Error(logger).Log("err", err)
		return nil, ErrCmdRepository
	}

	return newUser.(*model.User), nil
}