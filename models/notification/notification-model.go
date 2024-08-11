package notification

import (
	rpBase "github.com/IT-RushCode/rush_pkg/models"
)

// Notification модель
type Notification struct {
	ID        uint   `gorm:"primaryKey"`
	Type      string `gorm:"type:varchar(50);not null"`
	Recipient string `gorm:"type:varchar(255);not null"`
	Message   string `gorm:"type:text;not null"`

	rpBase.BaseModel
}

type Notifications []Notification

func (Notification) TableName() string {
	return "Notifications"
}
