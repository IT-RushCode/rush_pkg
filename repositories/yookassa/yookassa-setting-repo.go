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

// Полное обновление
func (r *yookassasettingRepository) SaveByPointID(ctx context.Context, data *models.YooKassaSetting) (*models.YooKassaSetting, error) {
	if err := r.db.WithContext(ctx).
		Where("point_id = ?", data.PointID).
		Save(&data).
		Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).
		Where("point_id = ?", data.PointID).
		First(&data).
		Error; err != nil {
		return nil, err
	}

	return data, nil
}
