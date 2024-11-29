package models

import (
	"gorm.io/gorm"
)

type NotificationDevice struct {
	ID                   uint   `gorm:"primaryKey"`
	UserID               *uint  `gorm:"index"`
	DeviceToken          string `gorm:"unique"`
	NotificationsEnabled *bool  `gorm:"default:true"`
	IsAuthenticated      *bool  `gorm:"default:false"`

	BaseModel
}

// Настройки Notification
type NotificationDevices []NotificationDevice

func (NotificationDevice) TableName() string {
	return "NotificationDevices"
}

func (m *NotificationDevice) BeforeCreate(db *gorm.DB) (err error) {
	if err := CheckSequence(m.TableName(), db); err != nil {
		return err
	}
	return nil
}
