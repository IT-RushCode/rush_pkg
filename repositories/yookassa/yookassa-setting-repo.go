package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/IT-RushCode/rush_pkg/models"
	rp "github.com/IT-RushCode/rush_pkg/repositories/base"
	"github.com/IT-RushCode/rush_pkg/utils"

	"gorm.io/gorm"
)

type YooKassaSettingRepository interface {
	rp.BaseRepository
	UpdateByPointID(ctx context.Context, data *models.YooKassaSetting) (*models.YooKassaSetting, error)
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
func (r *yookassasettingRepository) UpdateByPointID(ctx context.Context, data *models.YooKassaSetting) (*models.YooKassaSetting, error) {
	if err := r.db.WithContext(ctx).
		Where("point_id = ?", data.PointID).
		Updates(&data).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Print(err)
			return nil, utils.ErrRecordNotFound
		}
		if err := utils.HandleDuplicateKeyError(err); err != nil {
			log.Print(err)
			return nil, err
		}
		log.Print(err)
		return nil, utils.ErrInternal
	}

	if err := r.db.WithContext(ctx).
		Where("point_id = ?", data.PointID).
		First(&data).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Print(err)
			return nil, utils.ErrRecordNotFound
		}
		log.Print(err)
		return nil, utils.ErrInternal
	}

	return data, nil
}
