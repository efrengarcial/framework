package repository

import (
	"context"
	"github.com/efrengarcial/framework/internal/domain"

	"github.com/efrengarcial/framework/internal/platform/repository"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sagikazarmark/go-gin-gorm-opencensus/pkg/ocgorm"
)

type userGormRepository struct {
	repository.GormRepository
}
//https://stackoverflow.com/questions/32533992/golang-inheritance-using-setting-protected-values
//The final piece of advice would be to not return interfaces from NewX methods. You should (almost) always return a pointer to the struct.
// Generally you want to construct a type, and pass it into another method as an interface, not receive an interface in the first place.
func NewUserGormRepository(db *gorm.DB) *userGormRepository {

	repo := repository.NewGormRepository(db)
	return &userGormRepository{repo}
}

func (repo *userGormRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	user := domain.User{}
	orm := ocgorm.WithContext(ctx, repo.DB())
	err := orm.Preload("Authorities").Where("email = ?", email).First(&user).Error
	if err !=nil {
		if gorm.IsRecordNotFoundError(err) {
			return domain.User{}, domain.ErrAuthenticationFailure
		} else {
			return domain.User{}, errors.WithStack(err)
		}
	}
	return user, nil
}

func (repo *userGormRepository) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	user := domain.User{}
	orm := ocgorm.WithContext(ctx, repo.DB())
	err := orm.Preload("Authorities").Where("login = ?", login).First(&user).Error
	if err !=nil {
		if gorm.IsRecordNotFoundError(err) {
			return domain.User{}, domain.ErrAuthenticationFailure
		} else {
			return domain.User{}, errors.WithStack(err)
		}
	}

	return user, nil
}

func (repo *userGormRepository) FindOneByLogin(ctx context.Context, login string) (domain.User, error) {
	user := domain.User{}
	orm := ocgorm.WithContext(ctx, repo.DB())
	var err error
	err = orm.Where("login = ?", login).First(&user).Error

	if gorm.IsRecordNotFoundError(err) {
		return domain.User{}, nil
	}

	if err != nil {
		return  domain.User{},  errors.WithStack(err)
	}

	return user, nil
}

func (repo *userGormRepository) FindOneByEmail(ctx context.Context, login string) (domain.User, error) {
	user := domain.User{}
	orm := ocgorm.WithContext(ctx, repo.DB())
	var err error
	err = orm.Where("email = ?", login).First(&user).Error

	if gorm.IsRecordNotFoundError(err) {
		return domain.User{}, nil
	}

	if err != nil {
		return  domain.User{}, errors.WithStack(err)
	}

	return user, nil
}
