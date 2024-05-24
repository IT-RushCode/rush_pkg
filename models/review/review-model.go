package review

import (
	"time"

	"gorm.io/gorm"
)

// Отзывы
type Review struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	PointID   uint
	Comment   string `gorm:"type:text"`
	Rating    int
	Status    bool           `gorm:"default:true"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Reviews []Review

func (Review) TableName() string {
	return "Reviews"
}
