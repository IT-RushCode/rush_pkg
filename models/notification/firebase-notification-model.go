package notification

import (
	"time"

	rpBase "github.com/IT-RushCode/rush_pkg/models"
)

// FirebaseNotification модель
type FirebaseNotification struct {
	ID          uint   `gorm:"primaryKey"`
	DeviceToken string `gorm:"type:varchar(255);not null"`
	Title       string `gorm:"type:varchar(255);not null"`
	Body        string `gorm:"type:text;not null"`
	Data        string `gorm:"type:jsonb"`
	Status      string `gorm:"type:varchar(50);not null;default:'pending'"`
	SentAt      *time.Time
	FailedAt    *time.Time
	Error       string `gorm:"type:text"`
	rpBase.BaseModel
}

type FirebaseNotifications []FirebaseNotification

func (FirebaseNotification) TableName() string {
	return "FirebaseNotifications"
}
