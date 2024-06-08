package review

import (
	rpBase "github.com/IT-RushCode/rush_pkg/models"
	rpAuth "github.com/IT-RushCode/rush_pkg/models/auth"
)

// Отзывы
type Review struct {
	ID      uint `gorm:"primaryKey"`
	UserID  uint
	Comment string `gorm:"type:text"`
	Rating  int
	Status  *bool
	User    rpAuth.User `gorm:"foreignKey:UserID"`
	rpBase.BaseModel
}

type Reviews []Review

func (Review) TableName() string {
	return "Reviews"
}
