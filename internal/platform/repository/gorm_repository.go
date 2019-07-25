package repository

import (
	"github.com/efrengarcial/framework/internal/platform/database"
	"github.com/efrengarcial/framework/internal/platform/service"
	"github.com/jinzhu/gorm"
)

type GormRepository struct {
	DB *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return GormRepository{DB:db}
}

func (gr GormRepository) Insert(model service.IModel) (service.IModel, error){
	if err := model.Validate(); err != nil{
		return nil, err
	}
	if err := gr.DB.Create(model).Error; err != nil{
		return nil, err
	}
	return model, nil
}

func (gr GormRepository) Update(model service.IModel) error {
	if err := model.Validate(); err != nil{
		return err
	}
	return gr.DB.Save(model).Error
}

func (gr GormRepository) Save(model service.IModel) (uint64, error){
	if err := model.Validate(); err != nil{
		return 0, err
	}
	if err := gr.DB.Save(model).Error; err != nil{
		return 0, err
	}
	return model.GetID(), nil
}

func (gr GormRepository) Find(receiver service.IModel, id uint64) error {
	return gr.DB.Where("id = ?", id).Find(receiver).Error
}

func (gr GormRepository) FindFirst(receiver service.IModel, where string, args ...interface{}) error {
	return gr.DB.Where(where, args...).Limit(1).Find(receiver).Error
}

func (gr GormRepository) FindAll(result interface{}, where string, args ...interface{}) (err error){
	err = gr.DB.Where(where, args...).Find(result).Error
	return
}

func (gr GormRepository) FindAllPageable(pageable *service.Pageable, result interface{},  where string, args ...interface{} ) (*database.Pagination, error) {
	//http://jinzhu.me/gorm/crud.html#query
	//err := gr.DB.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&result)
	//gr.DB = gr.DB.Model(pageable.Model).Where(where, args)
	p := &database.Param{
		DB:      gr.DB.Where(where, args),
		Page:    pageable.Page,
		Limit:   pageable.Limit,
		OrderBy: pageable.OrderBy,
		ShowSQL: pageable.ShowSQL,
	}
	return database.Pagging(p, result)
}

func (gr GormRepository) Delete(model service.IModel, where string, args ...interface{}) error {
	return gr.DB.Where(where, args...).Delete(&model).Error
}

func (gr GormRepository) NewRecord(model service.IModel) bool {
	return gr.DB.NewRecord(&model)
}
