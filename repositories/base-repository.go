package repositories

import (
	"gorm.io/gorm"
)

// Repository интерфейс представляет базовый набор методов для работы с сущностями
type BaseRepository interface {
	Create(data interface{}) error
	FindByID(id uint, data interface{}) error
	Update(data interface{}) error
	Delete(data interface{}) error
	UpdateField(id uint, field string, value interface{}, data interface{}) error
	SoftDelete(data interface{}) error
	Filter(filters map[string]interface{}, entities interface{}) error
}

// BaseRepository представляет базовую структуру для репозиториев
type baseRepository struct {
	db *gorm.DB
}

// NewBaseRepository создает новый экземпляр базового репозитория
func NewBaseRepository(db *gorm.DB) BaseRepository {
	return &baseRepository{db: db}
}

// Реализация методов интерфейса Repository
func (r *baseRepository) Create(data interface{}) error {
	return r.db.Create(data).Error
}

func (r *baseRepository) FindByID(id uint, data interface{}) error {
	return r.db.First(data, id).Error
}

func (r *baseRepository) Update(data interface{}) error {
	return r.db.Save(data).Error
}

func (r *baseRepository) Delete(data interface{}) error {
	return r.db.Delete(data).Error
}

func (r *baseRepository) UpdateField(id uint, field string, value interface{}, data interface{}) error {
	return r.db.Model(data).Where("id = ?", id).Update(field, value).Error
}

func (r *baseRepository) SoftDelete(data interface{}) error {
	return r.db.Model(data).Update("deleted_at", gorm.DeletedAt{}).Error
}

func (r *baseRepository) Filter(filters map[string]interface{}, entities interface{}) error {
	query := r.db
	for key, value := range filters {
		query = query.Where(key, value)
	}
	return query.Find(entities).Error
}
