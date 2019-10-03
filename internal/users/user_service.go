package users

import (
	"context"
	"github.com/efrengarcial/framework/internal/domain"
	"strings"
	"time"

	"github.com/efrengarcial/framework/internal/platform"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)


// UserService describes the service.
type UserService interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	FindAll(ctx context.Context, pageable *domain.Pageable, result interface{},  where string, args ...interface{})(*domain.Pagination, error)
}


// service implements the User Service
type userService struct {
	repository UserRepository
	logger     *logrus.Logger
}

// NewService creates and returns a new User service instance
func NewService(rep UserRepository, logger *logrus.Logger) *userService {
	return &userService {
		repository: rep,
		logger:     logger,
	}
}

//https://github.com/pkg/errors
//https://medium.com/@hussachai/error-handling-in-go-a-quick-opinionated-guide-9199dd7c7f76
func (service *userService) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	var  (
		err error
		existingUser *domain.User
	)

	if user.ID > 0 {
		return nil, ErrIdExist
	}

	if  existingUser  , err = service.repository.FindOneByLogin(ctx, strings.ToLower(user.Login)); existingUser != nil {
		return nil, ErrLoginAlreadyUsed
	}
	if err != nil { return nil, err }

	if  existingUser  , err =  service.repository.FindOneByEmail(ctx, strings.ToLower(user.Email));  existingUser != nil {
		return nil, ErrEmailAlreadyUsed
	}
	if err != nil { return nil, err}

	if len(user.LangKey) ==0 {
		user.LangKey = "en"
	}

	// Generates a hashed version of our password
	//randomPassword, _ := platform.GeneratePassword()
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil { return nil, err}

	user.Password = string(hashedPass)
	resetKey, _ := platform.GenerateResetKey()
	user.ResetKey = resetKey
	user.ResetDate = time.Now()
	user.Activated = true

	newUser, err  := service.repository.Insert(ctx, user)
	if err != nil { return nil, err}


	return newUser.(*domain.User), nil
}


func (service *userService) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	var  (
		err error
		existingUser *domain.User
	)

	if  existingUser  , err =  service.repository.FindOneByEmail(ctx, strings.ToLower(user.Email)); existingUser != nil && user.ID !=  existingUser.ID {
		return nil, ErrEmailAlreadyUsed
	}
	if err != nil { return nil, err }

	if  existingUser  , err =  service.repository.FindOneByLogin(ctx, strings.ToLower(user.Login)); existingUser != nil && user.ID !=  existingUser.ID {
		return nil, ErrLoginAlreadyUsed
	}
	if err != nil { return nil, err }

	err  = service.repository.Update(user)
	if err != nil { return nil, err }

	return user, nil
}


func (service *userService) FindAll(ctx context.Context, pageable *domain.Pageable, result interface{}, where string, args ...interface{}) (*domain.Pagination, error){
	return service.repository.FindAllPageable(ctx, pageable, result, where, args...)
}
