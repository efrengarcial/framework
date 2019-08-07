package repository

import (
	"context"
	"github.com/efrengarcial/framework/internal/platform/repository"
	"github.com/efrengarcial/framework/internal/users/service"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sagikazarmark/go-gin-gorm-opencensus/pkg/ocgorm"
)

type userGormRepository struct {
	repository.GormRepository
}

func NewUserGormRepository(db *gorm.DB) service.UserRepository {

	repo := repository.GormRepository{DB: db}
	return &userGormRepository{repo}
}

func (repo *userGormRepository) GetByEmail(email string) (*service.User, error) {
	user := &service.User{}
	if err := repo.DB.Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return user, nil
}

func (repo *userGormRepository) GetByLogin(login string) (*service.User, error) {
	user := &service.User{}
	if err := repo.DB.Where("login = ?", login).
		First(&user).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return user, nil
}

func (repo *userGormRepository) FindOneByLogin(ctx context.Context, login string) (*service.User, error) {
	user := &service.User{}
	orm := ocgorm.WithContext(ctx, repo.DB)
	var err error
	err = orm.Where("login = ?", login).First(&user).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return  nil,  errors.WithStack(err)
	}

	return user, nil
}

func (repo *userGormRepository) FindOneByEmail(ctx context.Context, login string) (*service.User, error) {
	user := &service.User{}
	orm := ocgorm.WithContext(ctx, repo.DB)
	var err error
	err = orm.Where("email = ?", login).First(&user).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return  nil, errors.WithStack(err)
	}

	return user, nil
}
