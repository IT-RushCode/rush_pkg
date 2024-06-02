package review

import rpBase "github.com/IT-RushCode/rush_pkg/models"

// Отзывы
type Review struct {
	ID      uint `gorm:"primaryKey"`
	UserID  uint
	PointID uint
	Comment string `gorm:"type:text"`
	Rating  int
	Status  *bool
	rpBase.BaseModel
}

type Reviews []Review

func (Review) TableName() string {
	return "Reviews"
}
