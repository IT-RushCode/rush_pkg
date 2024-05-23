package repositories

import (
	"context"

	"gorm.io/gorm"
)

// Repository интерфейс представляет базовый набор методов для работы с сущностями
type BaseRepository interface {
	GetAll(ctx context.Context, offset, limit uint, data interface{}) (int64, error)
	Create(ctx context.Context, data interface{}) (interface{}, error)
	FindByID(ctx context.Context, id uint, data interface{}) error
	Update(ctx context.Context, data interface{}) (interface{}, error)
	Delete(ctx context.Context, data interface{}) error
	UpdateField(ctx context.Context, id uint, field string, value interface{}, data interface{}) (interface{}, error)
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

// Получение всех или с пагинацией
func (r *baseRepository) GetAll(ctx context.Context, offset, limit uint, data interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(data)

	// Получить общее количество записей
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	// Применить пагинацию, если необходимо
	if limit > 0 || offset > 0 {
		query = query.Offset(int(offset)).Limit(int(limit))
	}

	// Получить данные с учетом пагинации
	if err := query.Find(data).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// Создание записи
func (r *baseRepository) Create(ctx context.Context, data interface{}) (interface{}, error) {
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Поиск записи по ID
func (r *baseRepository) FindByID(ctx context.Context, id uint, data interface{}) error {
	return r.db.WithContext(ctx).
		First(data, id).Error
}

// Полное обновление
func (r *baseRepository) Update(ctx context.Context, data interface{}) (interface{}, error) {
	if err := r.db.WithContext(ctx).Save(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Частичное обновление
func (r *baseRepository) UpdateField(
	ctx context.Context,
	id uint, field string, value interface{}, data interface{},
) (interface{}, error) {
	if err := r.db.WithContext(ctx).
		Model(data).
		Where("id = ?", id).
		Update(field, value).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).First(data, id).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Перманентное удаление
func (r *baseRepository) Delete(ctx context.Context, data interface{}) error {
	return r.db.WithContext(ctx).
		Delete(data).Error
}

// Мягкое удаление (устанавливается время deleted_at для записи)
func (r *baseRepository) SoftDelete(ctx context.Context, data interface{}) error {
	return r.db.WithContext(ctx).
		Model(data).
		Update("deleted_at", gorm.DeletedAt{}).Error
}

// Получение данных по фильтру
func (r *baseRepository) Filter(
	ctx context.Context,
	filters map[string]interface{}, entities interface{},
) error {
	query := r.db.WithContext(ctx)
	for key, value := range filters {
		query = query.Where(key, value)
	}
	return query.Find(entities).Error
}
