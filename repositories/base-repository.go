package repositories

import (
	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

func (r *BaseRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *BaseRepository[T]) FindByID(id uint) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	return &entity, err
}

func (r *BaseRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

func (r *BaseRepository[T]) Delete(entity *T) error {
	return r.db.Delete(entity).Error
}

func (r *BaseRepository[T]) UpdateField(id uint, field string, value interface{}) error {
	return r.db.Model(new(T)).Where("id = ?", id).Update(field, value).Error
}

func (r *BaseRepository[T]) SoftDelete(entity *T) error {
	return r.db.Model(entity).Update("deleted_at", gorm.DeletedAt{}).Error
}

func (r *BaseRepository[T]) Filter(filters map[string]interface{}) ([]T, error) {
	var entities []T
	query := r.db
	for key, value := range filters {
		query = query.Where(key, value)
	}
	err := query.Find(&entities).Error
	return entities, err
}
