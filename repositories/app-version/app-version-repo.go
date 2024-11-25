package repositories

import (
	"context"
	"errors"

	"github.com/IT-RushCode/rush_pkg/models"
	"gorm.io/gorm"
)

type AppVersionRepository struct {
	db *gorm.DB
}

func NewAppVersionRepository(db *gorm.DB) *AppVersionRepository {
	return &AppVersionRepository{db: db}
}

func (r *AppVersionRepository) GetLatest(ctx context.Context) (*models.AppVersion, error) {
	var version models.AppVersion
	err := r.db.WithContext(ctx).Order("created_at desc").First(&version).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &version, err
}

func (r *AppVersionRepository) Create(ctx context.Context, version *models.AppVersion) error {
	return r.db.WithContext(ctx).Create(version).Error
}

func (r *AppVersionRepository) Update(ctx context.Context, version *models.AppVersion) error {
	return r.db.WithContext(ctx).Save(version).Error
}
