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

type yooKassaSettingRepository struct {
	rp.BaseRepository
	db *gorm.DB
}

func NewYooKassaSettingRepository(db *gorm.DB) YooKassaSettingRepository {
	return &yooKassaSettingRepository{
		BaseRepository: rp.NewBaseRepository(db),
		db:             db,
	}
}

// Полное обновление по point_id
func (r *yooKassaSettingRepository) SaveByPointID(ctx context.Context, data *models.YooKassaSetting) (*models.YooKassaSetting, error) {
	var existingRecord models.YooKassaSetting

	// Используем FirstOrCreate для поиска или создания записи
	if err := r.db.WithContext(ctx).
		Where("point_id = ?", data.PointID).
		Attrs(models.YooKassaSetting{
			StoreID:   data.StoreID,
			SecretKey: data.SecretKey,
			IsTest:    data.IsTest,
			Status:    data.Status,
		}).
		FirstOrCreate(&existingRecord).Error; err != nil {
		return nil, err
	}

	// Если запись существовала, обновляем поля
	if existingRecord.ID != 0 {
		if err := r.db.WithContext(ctx).Model(&existingRecord).Updates(data).Error; err != nil {
			return nil, err
		}
	}

	return &existingRecord, nil
}
