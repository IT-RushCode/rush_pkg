package repositories

import (
	rp "github.com/IT-RushCode/rush_pkg/repositories/base"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	rp.BaseRepository
}

type permissionRepository struct {
	db *gorm.DB
	rp.BaseRepository
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{
		BaseRepository: rp.NewBaseRepository(db),
		db:             db,
	}
}
