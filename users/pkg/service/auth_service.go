package service


import (
	"context"
	"github.com/efrengarcial/framework/users/pkg/model"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// AuthService has the logic authentication
type AuthService interface {
	Auth(ctx context.Context, in *model.User) (*model.Token, error)
	ValidateToken(ctx context.Context, in *model.Token) (*model.Token, error)
}

type authService struct {
	repo         Repository
	tokenService TokenService
}

func (srv *authService) Auth(ctx context.Context, req *model.User, res *model.Token) error {
	log.Println("Logging in with:", req.Email, req.Password)
	user, err := srv.repo.GetByEmail(req.Email)
	log.Println(user, err)
	if err != nil {
		return err
	}

	// Compares our given password against the hashed password
	// stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}

	token, err := srv.tokenService.Encode(user)
	if err != nil {
		return err
	}
	res.Token = token
	return nil
}

func (srv *authService) ValidateToken(ctx context.Context, req *model.Token, res *model.Token) error {

	// Decode token
	claims, err := srv.tokenService.Decode(req.Token)

	if err != nil {
		return err
	}

	if claims.User.ID == 0 {
		return errors.New("invalid user")
	}

	res.Valid = true

	return nil
}