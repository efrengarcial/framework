package repository

import (
	"context"
	"github.com/efrengarcial/framework/internal/platform/database"
	"github.com/efrengarcial/framework/internal/platform/service"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sagikazarmark/go-gin-gorm-opencensus/pkg/ocgorm"
)

type GormRepository struct {
	DB *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return GormRepository{DB:db}
}

func (gr GormRepository) Insert(ctx context.Context, model service.IModel) (service.IModel, error){
	orm := ocgorm.WithContext(ctx, gr.DB)
	if err := model.Validate(); err != nil{
		return nil,  errors.WithStack(err)
	}
	if err := orm.Create(model).Error; err != nil{
		return nil,  errors.WithStack(err)
	}
	return model, nil
}

func (gr GormRepository) Update(model service.IModel) error {
	if err := model.Validate(); err != nil{
		return  errors.WithStack(err)
	}
	if err := gr.DB.Save(model).Error; err != nil{
		return errors.WithStack(err)
	}
	return nil
}

func (gr GormRepository) Save(model service.IModel) (uint64, error){
	if err := model.Validate(); err != nil{
		return 0, errors.WithStack(err)
	}
	if err := gr.DB.Save(model).Error; err != nil{
		return 0, errors.WithStack(err)
	}
	return model.GetID(), nil
}

func (gr GormRepository) Find(receiver service.IModel, id uint64) error {
	if err := gr.DB.Where("id = ?", id).Find(receiver).Error; err != nil{
		return errors.WithStack(err)
	}
	return nil
}

func (gr GormRepository) FindFirst(receiver service.IModel, where string, args ...interface{}) error {
	if err := gr.DB.Where(where, args...).Limit(1).Find(receiver).Error; err != nil{
		return errors.WithStack(err)
	}
	return nil
}

func (gr GormRepository) FindAll(result interface{}, where string, args ...interface{}) (err error){
	if err := gr.DB.Where(where, args...).Find(result).Error; err != nil{
		return errors.WithStack(err)
	}
	return nil
}

func (gr GormRepository) FindAllPageable(ctx context.Context, pageable *service.Pageable, result interface{},  where string, args ...interface{} ) (*database.Pagination, error) {
	//http://jinzhu.me/gorm/crud.html#query
	//err := gr.DB.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&result)
	//gr.DB = gr.DB.Model(pageable.Model).Where(where, args)
	orm := ocgorm.WithContext(ctx, gr.DB)
	p := &database.Param{
		DB:      orm.Where(where, args),
		Page:    pageable.Page,
		Limit:   pageable.Limit,
		OrderBy: pageable.OrderBy,
		ShowSQL: pageable.ShowSQL,
	}
	pagination, err := database.Pagging(p, result)
	return pagination,  errors.WithStack(err)
}

func (gr GormRepository) Delete(model service.IModel, where string, args ...interface{}) error {
	if err :=gr.DB.Where(where, args...).Delete(&model).Error; err != nil{
		return errors.WithStack(err)
	}
	return nil
}

func (gr GormRepository) NewRecord(model service.IModel) bool {
	return gr.DB.NewRecord(&model)
}
