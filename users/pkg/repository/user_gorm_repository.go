package repository

import (
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/efrengarcial/framework/users/pkg/service"
)

type UserGormRepository struct {
	GormRepository
}

func NewUserGormRepository(repository GormRepository) service.UserRepository {
	return UserGormRepository{repository}
}

func (repo *UserGormRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := repo.DB.Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

