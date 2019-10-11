package service

import (
	"context"
	"github.com/efrengarcial/framework/internal/domain"
	"strings"
	"time"

	"github.com/efrengarcial/framework/internal/platform"
	"github.com/efrengarcial/framework/internal/user"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// service implements the User Service
type userService struct {
	repository  user.Repository
	logger     *logrus.Logger
}

// NewService creates and returns a new User service instance
func NewService(rep user.Repository, logger *logrus.Logger) *userService {
	return &userService {
		repository: rep,
		logger:     logger,
	}
}

//https://github.com/pkg/errors
//https://medium.com/@hussachai/error-handling-in-go-a-quick-opinionated-guide-9199dd7c7f76
func (service *userService) Create(ctx context.Context, us *domain.User) error {
	var  (
		err error
		existingUser  domain.User
	)

	if us.ID > 0 {
		return user.ErrIdExist
	}

	if  existingUser  , err = service.repository.FindOneByLogin(ctx, strings.ToLower(us.Login)); existingUser.ID != 0 {
		return user.ErrLoginAlreadyUsed
	}
	if err != nil { return err }

	if  existingUser  , err =  service.repository.FindOneByEmail(ctx, strings.ToLower(us.Email));  existingUser.ID != 0 {
		return user.ErrEmailAlreadyUsed
	}
	if err != nil { return err}

	if len(us.LangKey) ==0 {
		us.LangKey = "en"
	}

	// Generates a hashed version of our password
	//randomPassword, _ := platform.GeneratePassword()
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(us.Password), bcrypt.DefaultCost)
	if err != nil { return err}

	us.Password = domain.Password(hashedPass)
	resetKey, _ := platform.GenerateResetKey()
	us.ResetKey = resetKey
	us.ResetDate = time.Now()
	us.Activated = true

	err  = service.repository.Insert(ctx, us)
	if err != nil { return err}

	return nil
}


func (service *userService) Update(ctx context.Context, us *domain.User)  error {
	var  (
		err error
		existingUser domain.User
	)

	if  existingUser  , err =  service.repository.FindOneByEmail(ctx, strings.ToLower(us.Email));  us.ID !=  existingUser.ID {
		return user.ErrEmailAlreadyUsed
	}
	if err != nil { return err }

	if  existingUser  , err =  service.repository.FindOneByLogin(ctx, strings.ToLower(us.Login)); us.ID !=  existingUser.ID {
		return user.ErrLoginAlreadyUsed
	}
	if err != nil { return err }

	return service.repository.Update(us)
}


func (service *userService) FindAll(ctx context.Context, pageable *domain.Pageable, result interface{}, where string, args ...interface{}) (*domain.Pagination, error){
	return service.repository.FindAllPageable(ctx, pageable, result, where, args...)
}
