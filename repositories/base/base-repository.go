package base_repository

import (
	"context"
	"errors"

	"github.com/IT-RushCode/rush_pkg/utils"
	"gorm.io/gorm"
)

// BaseRepository интерфейс представляет базовый набор методов для работы с сущностями
type BaseRepository interface {
	GetAll(ctx context.Context, offset, limit uint, data interface{}) (int64, error)
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

// ----------- Реализация базовых методов интерфейса Repository -----------

// Получение всех или с пагинацией
func (r *baseRepository) GetAll(ctx context.Context, offset, limit uint, data interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(data)

	// Получить общее количество записей
	if err := query.Count(&count).Error; err != nil {
		return 0, utils.ErrInternal
	}

	// Применить пагинацию, если необходимо
	if limit > 0 || offset > 0 {
		query = query.Scopes(utils.Paginate(offset, limit))
	}

	// Получить данные с учетом пагинации
	if err := query.Find(data).Error; err != nil {
		return 0, utils.ErrInternal
	}

	return count, nil
}

// Создание записи
func (r *baseRepository) Create(ctx context.Context, data interface{}) error {
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		if err := utils.HandleDuplicateKeyError(err); err != nil {
			return err
		}
		return utils.ErrInternal
	}
	return nil
}

// Поиск записи по ID
func (r *baseRepository) FindByID(ctx context.Context, id uint, data interface{}) error {
	err := r.db.WithContext(ctx).First(data, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrRecordNotFound
		}
		return utils.ErrInternal
	}
	return nil
}

// Полное обновление
func (r *baseRepository) Update(ctx context.Context, data interface{}) error {
	if err := r.db.WithContext(ctx).
		Updates(data).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrRecordNotFound
		}
		if err := utils.HandleDuplicateKeyError(err); err != nil {
			return err
		}
		return utils.ErrInternal
	}

	if err := r.db.WithContext(ctx).
		First(data).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrRecordNotFound
		}
		return utils.ErrInternal
	}

	return nil
}

// Частичное обновление
func (r *baseRepository) UpdateField(ctx context.Context, id uint, field string, value interface{}, data interface{}) error {
	err := r.db.WithContext(ctx).
		Model(data).
		Where("id = ?", id).
		Update(field, value).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrRecordNotFound
		}
		return utils.ErrInternal
	}
	if err := r.db.WithContext(ctx).First(data, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrRecordNotFound
		}
		return utils.ErrInternal
	}
	return nil
}

// Перманентное удаление
func (r *baseRepository) Delete(ctx context.Context, data interface{}) error {
	err := r.db.WithContext(ctx).Delete(data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrRecordNotFound
		}
		return utils.ErrInternal
	}
	return nil
}

// Мягкое удаление (устанавливается время deleted_at для записи)
func (r *baseRepository) SoftDelete(ctx context.Context, data interface{}) error {
	err := r.db.WithContext(ctx).
		Model(data).
		Update("deleted_at", gorm.DeletedAt{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrRecordNotFound
		}
		return utils.ErrInternal
	}
	return nil
}

// Получение данных по фильтру
func (r *baseRepository) Filter(ctx context.Context, filters map[string]interface{}, entities interface{}) error {
	query := r.db.WithContext(ctx)
	for key, value := range filters {
		query = query.Where(key, value)
	}
	if err := query.Find(entities).Error; err != nil {
		return utils.ErrInternal
	}
	return nil
}
