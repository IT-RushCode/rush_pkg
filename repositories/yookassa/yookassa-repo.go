package repositories

import (
	rp "github.com/IT-RushCode/rush_pkg/repositories/base"

	"gorm.io/gorm"
)

type YooKassaSettingRepository interface {
	rp.BaseRepository
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
