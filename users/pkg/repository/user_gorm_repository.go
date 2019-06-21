package repository

import (
	"github.com/efrengarcial/framework/users/pkg/model"
	. "github.com/efrengarcial/framework/users/pkg/service"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type userGormRepository struct {
	GormRepository
}

func NewUserGormRepository(db *gorm.DB) UserRepository {

	repo := GormRepository{db}
	return &userGormRepository{repo}
}

func (repo *userGormRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := repo.DB.Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return user, nil
}

func (repo *userGormRepository) GetByLogin(login string) (*model.User, error) {
	user := &model.User{}
	if err := repo.DB.Where("login = ?", login).
		First(&user).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return user, nil
}

func (repo *userGormRepository) FindOneByLogin(login string) (*model.User, error) {
	user := &model.User{}
	var err error
	err = repo.DB.Where("login = ?", login).First(&user).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return  nil,  errors.WithStack(err)
	}

	return user, nil
}

func (repo *userGormRepository) FindOneByEmail(login string) (*model.User, error) {
	user := &model.User{}
	var err error
	err = repo.DB.Where("email = ?", login).First(&user).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return  nil, errors.WithStack(err)
	}

	return user, nil
}
