package base_repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/IT-RushCode/rush_pkg/dto"
	"github.com/IT-RushCode/rush_pkg/utils"
	"gorm.io/gorm"
)

// BaseRepository интерфейс представляет базовый набор методов для работы с сущностями
type BaseRepository interface {
	GetAll(ctx context.Context, data interface{}, dto *dto.GetAllRequest, pagination bool, preloads ...string) (int64, error)
	Create(ctx context.Context, data interface{}) error
	FindByID(ctx context.Context, id uint, data interface{}, preloads ...string) error
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
func (r *baseRepository) GetAll(ctx context.Context, data interface{}, dto *dto.GetAllRequest, pagination bool, preloads ...string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(data)

	// Применить preloads
	for _, preload := range preloads {
		if preload != "" {
			query = query.Preload(preload)
		}
	}

	// Применить фильтры, если они есть
	if len(dto.Filters) > 0 {
		for field, value := range dto.Filters {
			if field != "" && value != "" {
				switch {
				case value == "true" || value == "false":
					// Булевое значение
					query = query.Where(fmt.Sprintf("%s = ?", field), value)
				case isRFC3339(value):
					// Дата и время в формате RFC3339
					query = query.Where(fmt.Sprintf("%s = ?", field), value)
				case isDate(value):
					// Дата (формат: YYYY-MM-DD)
					query = query.Where(fmt.Sprintf("%s = ?", field), value)
				case isNumber(value):
					// Числовое значение
					query = query.Where(fmt.Sprintf("%s ILIKE ?", field), "%"+value+"%")
				default:
					// Строковое значение
					query = query.Where(fmt.Sprintf("%s ILIKE ?", field), "%"+value+"%")
				}

			}
		}
	}

	// Получить общее количество записей
	if err := query.Count(&count).Error; err != nil {
		return 0, utils.ErrInternal
	}

	// Применить сортировку
	if dto.SortBy != "" {
		order := "asc"
		if dto.OrderBy == "desc" {
			order = "desc"
		}
		query = query.Order(fmt.Sprintf("%s %s", dto.SortBy, order))
	}

	if pagination {
		// Получить данные с учетом пагинации
		if err := query.Scopes(utils.Paginate(dto.Offset, dto.Limit)).
			Find(data).Error; err != nil {
			return 0, utils.ErrInternal
		}
	} else {
		// Получить данные без учета пагинации
		if err := query.Find(data).Error; err != nil {
			return 0, utils.ErrInternal
		}
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
func (r *baseRepository) FindByID(ctx context.Context, id uint, data interface{}, preloads ...string) error {
	query := r.db.WithContext(ctx).Model(data)

	// Применить preloads
	for _, preload := range preloads {
		if preload != "" {
			query = query.Preload(preload)
		}
	}

	// Получить первую совпавшую запись
	err := query.First(data, id).Error
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

// isDate проверяет, является ли строка датой в формате YYYY-MM-DD
func isDate(value string) bool {
	_, err := time.Parse("2006-01-02", value)
	return err == nil
}

// isRFC3339 проверяет, является ли строка датой и временем в формате RFC3339
func isRFC3339(value string) bool {
	_, err := time.Parse(time.RFC3339, value)
	return err == nil
}

// isNumber проверяет, является ли строка числом
func isNumber(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil
}
