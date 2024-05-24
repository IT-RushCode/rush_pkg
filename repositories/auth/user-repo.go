package repositories

import (
	rp "github.com/IT-RushCode/rush_pkg/repositories"

	"gorm.io/gorm"
)

type UserRepository interface {
	rp.BaseRepository
}

type userRepository struct {
	rp.BaseRepository
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: rp.NewBaseRepository(db),
	}
}
