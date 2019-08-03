package service

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type LoginVM struct {
	UserName    string `json:"username"`
	Password 	string `json:"password"`
	RememberMe 	bool `json:"rememberMe"`
}


// AuthService has the logic authentication
type AuthService interface {
	Auth(ctx context.Context, req *LoginVM, res *Token) error
	ValidateToken(ctx context.Context, req *Token, res *Token) error
}

type authService struct {
	repo         UserRepository
	tokenService TokenService
	logger       *logrus.Logger
}

// NewService creates and returns a new Auth service instance
func NewAuthService(rep UserRepository, token TokenService, logger *logrus.Logger) AuthService {
	return &authService {
		repo: rep,
		tokenService: token,
		logger:     logger,
	}
}

func (service *authService) Auth(ctx context.Context, req *LoginVM, res *Token) error {
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

func (service *authService) ValidateToken(ctx context.Context, req *Token, res *Token) error {

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
