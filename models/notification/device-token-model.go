package notification

import (
	rpBase "github.com/IT-RushCode/rush_pkg/models"
	rpAuth "github.com/IT-RushCode/rush_pkg/models/auth"
)

// DeviceToken модель для FirebaseNotification
type DeviceToken struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       uint   `gorm:"primaryKey;index;not null;autoIncrement:false"`
	DeviceToken  string `gorm:"type:varchar(255);not null"`
	PushIsActive string `gorm:"type:varchar(50);not null;"`

	User rpAuth.User `gorm:"foreignKey:UserID"`

	rpBase.BaseModel
}

type DeviceTokens []DeviceToken

func (DeviceToken) TableName() string {
	return "DeviceTokens"
}
