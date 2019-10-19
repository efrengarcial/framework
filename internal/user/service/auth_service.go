package service

import (
	"context"
	"github.com/efrengarcial/framework/internal/domain"
	"github.com/efrengarcial/framework/internal/user"
	"time"

	"github.com/efrengarcial/framework/internal/platform/auth"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo          user.Repository
	authenticator *auth.Authenticator
	logger        *logrus.Logger
}

// NewService creates and returns a new Auth service instance
func NewAuthService(rep user.Repository, authenticator *auth.Authenticator, logger *logrus.Logger) *authService {
	return &authService {
		repo: rep,
		authenticator: authenticator,
		logger:     logger,
	}
}

func (service *authService) Auth(ctx context.Context, req *domain.LoginVM, tkn *domain.Token) error {
	user, err := service.repo.GetByLogin(ctx , req.UserName)

	if err != nil {
		return err
	}

	// Compares our given password against the hashed password
	// stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return domain.ErrAuthenticationFailure
	}
	// If we are this far the request is valid. Create some claims for the user
	// and generate their token.
	claims := auth.NewClaims(user.Login, user.GetRoles(), time.Now(), time.Hour)

	tkn.Token, err = service.authenticator.GenerateToken(claims)
	if err != nil {
		return errors.Wrap(err, "generating token")
	}
	tkn.Valid = true
	return nil
}
