package repository

import (
	"context"

	"github.com/efrengarcial/framework/internal/platform/database"
	base "github.com/efrengarcial/framework/internal/platform/model"
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

func (repo GormRepository) Insert(ctx context.Context, model base.IModel) (base.IModel, error){
	orm := ocgorm.WithContext(ctx, repo.db)
	if err := model.Validate(); err != nil{
		return nil,  errors.WithStack(err)
	}
	if err := orm.Create(model).Error; err != nil{
		return nil,  errors.WithStack(err)
	}
	return model, nil
}

func (repo GormRepository) Update(model base.IModel) error {
	if err := model.Validate(); err != nil{
		return  errors.WithStack(err)
	}
	if err := repo.db.Save(model).Error; err != nil{
		return errors.WithStack(err)
	}
	return nil
}

func (repo GormRepository) Save(model base.IModel) (uint64, error){
	if err := model.Validate(); err != nil{
		return 0, errors.WithStack(err)
	}
	if err := repo.db.Save(model).Error; err != nil{
		return 0, errors.WithStack(err)
	}
	return model.GetID(), nil
}

func (repo GormRepository) Find(receiver base.IModel, id uint64) error {
	if err := repo.db.Where("id = ?", id).Find(receiver).Error; err != nil{
		return errors.WithStack(err)
	}
	return nil
}

func (repo GormRepository) FindFirst(receiver base.IModel, where string, args ...interface{}) error {
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

func (repo GormRepository) FindAllPageable(ctx context.Context, pageable *base.Pageable, result interface{},  where string, args ...interface{} ) (*base.Pagination, error) {
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

func (repo GormRepository) Delete(model base.IModel) error {
	if err := repo.db.Delete(&model).Error; err != nil{
		return errors.WithStack(err)
	}
	return nil
}

func (repo GormRepository) NewRecord(model base.IModel) bool {
	return repo.db.NewRecord(&model)
}
