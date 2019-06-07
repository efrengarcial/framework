package service

import (
	"context"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type LoginVM struct {
	UserName    string `json:"username"`
	Password 	string `json:"password"`
	RememberMe 	bool `json:"rememberMe"`
}


// AuthService has the logic authentication
type AuthService interface {
	Auth(ctx context.Context, req *LoginVM, res *model.Token) error
	ValidateToken(ctx context.Context, req *model.Token, res *model.Token) error
}

type authService struct {
	repo         UserRepository
	tokenService TokenService
	logger     log.Logger
}

// NewService creates and returns a new Auth service instance
func NewAuthService(rep UserRepository, token TokenService, logger log.Logger) AuthService {
	return &authService {
		repo: rep,
		tokenService: token,
		logger:     logger,
	}
}

func (service *authService) Auth(ctx context.Context, req *LoginVM, res *model.Token) error {
	user, err := service.repo.GetByLogin(req.UserName)
	if err != nil {
		return err
	}

	// Compares our given password against the hashed password
	// stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}

	token, err := service.tokenService.Encode(user)
	if err != nil {
		return err
	}
	res.Token = token
	return nil
}

func (service *authService) ValidateToken(ctx context.Context, req *model.Token, res *model.Token) error {

	// Decode token
	claims, err := service.tokenService.Decode(req.Token)

	if err != nil {
		return err
	}

	if claims.User.ID == 0 {
		return errors.New("invalid user")
	}

	res.Valid = true

	return nil
}
