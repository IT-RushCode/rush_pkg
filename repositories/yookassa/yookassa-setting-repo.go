package repositories

import (
	"context"

	"github.com/IT-RushCode/rush_pkg/models"
	rp "github.com/IT-RushCode/rush_pkg/repositories/base"

	"gorm.io/gorm"
)

type YooKassaSettingRepository interface {
	rp.BaseRepository
	SaveByPointID(ctx context.Context, data *models.YooKassaSetting) (*models.YooKassaSetting, error)
}

type yookassasettingRepository struct {
	rp.BaseRepository
	db *gorm.DB
}

func NewYooKassaSettingRepository(db *gorm.DB) YooKassaSettingRepository {
	return &yookassasettingRepository{
		BaseRepository: rp.NewBaseRepository(db),
		db:             db,
	}
}

// Полное обновление по point_id
func (r *yookassasettingRepository) SaveByPointID(ctx context.Context, data *models.YooKassaSetting) (*models.YooKassaSetting, error) {
	var existingRecord models.YooKassaSetting

	// Проверяем, существует ли запись с данным point_id
	err := r.db.WithContext(ctx).
		Where("point_id = ?", data.PointID).
		First(&existingRecord).
		Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Если запись существует, обновляем её
	if err == nil {
		if err := r.db.WithContext(ctx).
			Model(&existingRecord).
			Updates(data).
			Error; err != nil {
			return nil, err
		}
		return &existingRecord, nil
	}

	// Если запись не найдена, создаем новую
	if err == gorm.ErrRecordNotFound {
		if err := r.db.WithContext(ctx).
			Create(data).
			Error; err != nil {
			return nil, err
		}
	}

	return data, nil
}
