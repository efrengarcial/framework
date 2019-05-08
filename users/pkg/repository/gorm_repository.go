package repository

import (
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/efrengarcial/framework/users/pkg/service"
	"github.com/jinzhu/gorm"
)

type GormRepository struct {
	Db *gorm.DB
}

func NewGormRepository(db *gorm.DB) service.Repository {
	return GormRepository{Db:db}
}

func (gr GormRepository) Insert(model model.IModel) (model.IModel, error){
	if err := model.Validate(); err != nil{
		return nil, err
	}
	if err := gr.Db.Create(model).Error; err != nil{
		return nil, err
	}
	return model, nil
}

func (gr GormRepository) Update(model model.IModel) error {
	if err := model.Validate(); err != nil{
		return err
	}
	return gr.Db.Save(model).Error
}

func (gr GormRepository) Save(model model.IModel) (uint64, error){
	if err := model.Validate(); err != nil{
		return 0, err
	}
	if err := gr.Db.Save(model).Error; err != nil{
		return 0, err
	}
	return model.GetID(), nil
}

func (gr GormRepository) FindById(receiver model.IModel, id uint64) error {
	return gr.Db.First(receiver, id).Error
}

func (gr GormRepository) FindFirst(receiver model.IModel, where string, args ...interface{}) error {
	return gr.Db.Where(where, args...).Limit(1).Find(receiver).Error
}

func (gr GormRepository) FindAll(models interface{}, where string, args ...interface{}) (err error){
	err = gr.Db.Where(where, args...).Find(models).Error
	return
}

func (gr GormRepository) Delete(model model.IModel, where string, args ...interface{}) error {
	return gr.Db.Where(where, args...).Delete(&model).Error
}

func (gr GormRepository) NewRecord(model model.IModel) bool {
	return gr.Db.NewRecord(&model)
}
