package repositories

import (
	"context"

	"gorm.io/gorm"
)

// Repository интерфейс представляет базовый набор методов для работы с сущностями
type BaseRepository interface {
	Create(ctx context.Context, data interface{}) error
	FindByID(ctx context.Context, id uint, data interface{}) error
	Update(ctx context.Context, data interface{}) error
	Delete(ctx context.Context, data interface{}) error
	UpdateField(ctx context.Context, id uint, field string, value interface{}, data interface{}) error
	SoftDelete(ctx context.Context, data interface{}) error
	Filter(ctx context.Context, filters map[string]interface{}, entities interface{}) error
}

// BaseRepository представляет базовую структуру для репозиториев
type baseRepository struct {
	db *gorm.DB
}

// NewBaseRepository создает новый экземпляр базового репозитория
func NewBaseRepository(db *gorm.DB) BaseRepository {
	return &baseRepository{
		db: db,
	}
}

// ----------- Реализация методов интерфейса Repository -----------

func (r *baseRepository) GetAll(ctx context.Context, limit, offset uint, data interface{}) error {
	return r.db.WithContext(ctx).
		Find(data).Error
}

func (r *baseRepository) Create(ctx context.Context, data interface{}) error {
	return r.db.WithContext(ctx).
		Create(data).Error
}

func (r *baseRepository) FindByID(ctx context.Context, id uint, data interface{}) error {
	return r.db.WithContext(ctx).
		First(data, id).Error
}

func (r *baseRepository) Update(ctx context.Context, data interface{}) error {
	return r.db.WithContext(ctx).
		Save(data).Error
}

func (r *baseRepository) Delete(ctx context.Context, data interface{}) error {
	return r.db.WithContext(ctx).
		Delete(data).Error
}

func (r *baseRepository) UpdateField(ctx context.Context, id uint, field string, value interface{}, data interface{}) error {
	return r.db.WithContext(ctx).
		Model(data).
		Where("id = ?", id).
		Update(field, value).Error
}

func (r *baseRepository) SoftDelete(ctx context.Context, data interface{}) error {
	return r.db.WithContext(ctx).
		Model(data).
		Update("deleted_at", gorm.DeletedAt{}).Error
}

func (r *baseRepository) Filter(ctx context.Context, filters map[string]interface{}, entities interface{}) error {
	query := r.db.WithContext(ctx)
	for key, value := range filters {
		query = query.Where(key, value)
	}
	return query.Find(entities).Error
}
