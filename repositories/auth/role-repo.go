package repositories

import (
	rp "github.com/IT-RushCode/rush_pkg/repositories"

	"gorm.io/gorm"
)

type RoleRepository interface {
	rp.BaseRepository
}

type roleRepository struct {
	db *gorm.DB
	rp.BaseRepository
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		BaseRepository: rp.NewBaseRepository(db),
		db:             db,
	}
}
