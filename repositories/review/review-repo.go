package repositories

import (
	"gorm.io/gorm"

	rp "github.com/IT-RushCode/rush_pkg/repositories/base"
)

type ReviewRepository interface {
	rp.BaseRepository
}

type reviewRepository struct {
	rp.BaseRepository
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{
		BaseRepository: rp.NewBaseRepository(db),
	}
}
