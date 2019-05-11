package repository

import (
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/efrengarcial/framework/users/pkg/service"
	"github.com/jinzhu/gorm"
)

type GormRepository struct {
	DB *gorm.DB
}

func NewGormRepository(db *gorm.DB) service.Repository {
	return GormRepository{DB:db}
}

func (gr GormRepository) Insert(model model.IModel) (model.IModel, error){
	if err := model.Validate(); err != nil{
		return nil, err
	}
	if err := gr.DB.Debug().Create(model).Error; err != nil{
		return nil, err
	}
	return model, nil
}

func (gr GormRepository) Update(model model.IModel) error {
	if err := model.Validate(); err != nil{
		return err
	}
	return gr.DB.Save(model).Error
}

func (gr GormRepository) Save(model model.IModel) (uint64, error){
	if err := model.Validate(); err != nil{
		return 0, err
	}
	if err := gr.DB.Save(model).Error; err != nil{
		return 0, err
	}
	return model.GetID(), nil
}

func (gr GormRepository) Find(receiver model.IModel, id uint64) error {
	return gr.DB.Where("id = ?", id).Find(receiver).Error
}

func (gr GormRepository) FindFirst(receiver model.IModel, where string, args ...interface{}) error {
	return gr.DB.Where(where, args...).Limit(1).Find(receiver).Error
}

func (gr GormRepository) FindAll(result interface{}, where string, args ...interface{}) (err error){
	err = gr.DB.Where(where, args...).Find(result).Error
	return
}

func (gr GormRepository) FindAllPageable(pageable model.Pageable, result interface{},  where string, args ...interface{} ) *pagination.Paginator {
	gr.DB = gr.DB.Model(pageable.Model).Where(where, args)
	p := &pagination.Param{
		DB:      gr.DB,
		Page:    pageable.Page,
		Limit:   pageable.Limit,
		OrderBy: pageable.OrderBy,
	}
	return pagination.Paging(p, result)
}

func (gr GormRepository) Delete(model model.IModel, where string, args ...interface{}) error {
	return gr.DB.Where(where, args...).Delete(&model).Error
}

func (gr GormRepository) NewRecord(model model.IModel) bool {
	return gr.DB.NewRecord(&model)
}
