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
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) GormRepository {
	return GormRepository{db: db}
}

func (repo GormRepository) DB() *gorm.DB {
	return repo.db
}

func (repo GormRepository) Insert(ctx context.Context, model service.IModel) (service.IModel, error){
	orm := ocgorm.WithContext(ctx, repo.db)
	if err := model.Validate(); err != nil{
		return nil,  errors.WithStack(err)
	}
	if err := orm.Create(model).Error; err != nil{
		return nil,  errors.WithStack(err)
	}
	return model, nil
}

func (repo GormRepository) Update(model service.IModel) error {
	if err := model.Validate(); err != nil{
		return  errors.WithStack(err)
	}
	if err := repo.db.Save(model).Error; err != nil{
		return errors.WithStack(err)
	}
	return nil
}

func (repo GormRepository) Save(model service.IModel) (uint64, error){
	if err := model.Validate(); err != nil{
		return 0, errors.WithStack(err)
	}
	if err := repo.db.Save(model).Error; err != nil{
		return 0, errors.WithStack(err)
	}
	return model.GetID(), nil
}

func (repo GormRepository) Find(receiver service.IModel, id uint64) error {
	if err := repo.db.Where("id = ?", id).Find(receiver).Error; err != nil{
		return errors.WithStack(err)
	}
	return nil
}

func (repo GormRepository) FindFirst(receiver service.IModel, where string, args ...interface{}) error {
	if err := repo.db.Where(where, args...).Limit(1).Find(receiver).Error; err != nil{
		return errors.WithStack(err)
	}
	return nil
}

func (repo GormRepository) FindAll(result interface{}, where string, args ...interface{}) (err error){
	if err := repo.db.Where(where, args...).Find(result).Error; err != nil{
		return errors.WithStack(err)
	}
	return nil
}

func (repo GormRepository) FindAllPageable(ctx context.Context, pageable *service.Pageable, result interface{},  where string, args ...interface{} ) (*service.Pagination, error) {
	//http://jinzhu.me/gorm/crud.html#query
	//err := repo.db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&result)
	//repo.db = repo.db.Model(pageable.Model).Where(where, args)
	orm := ocgorm.WithContext(ctx, repo.db)
	p := &database.Param{
		DB:      orm.Where(where, args),
		Page:    pageable.Page,
		Limit:   pageable.Limit,
		OrderBy: pageable.OrderBy,
		ShowSQL: pageable.ShowSQL,
	}
	pagination, err := database.Pagging(p, result)
	if err !=nil {
		return nil,  errors.WithStack(err)
	}
	return pagination,  nil
}

func (repo GormRepository) Delete(model service.IModel) error {
	if err := repo.db.Delete(&model).Error; err != nil{
		return errors.WithStack(err)
	}
	return nil
}

func (repo GormRepository) NewRecord(model service.IModel) bool {
	return repo.db.NewRecord(&model)
}
