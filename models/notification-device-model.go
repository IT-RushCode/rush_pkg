package models

import (
	"gorm.io/gorm"
)

type NotificationDevice struct {
	ID                   uint   `gorm:"primaryKey;autoIncrement"`
	UserID               *uint  `gorm:"default:null"`      // ID пользователя, если пользователь авторизован
	DeviceToken          string `gorm:"type:varchar(255)"` // Токен устройства
	NotificationsEnabled *bool  `gorm:"default:true"`      // Статус активности получения уведомления пользователем

	BaseModel
}

// Настройки Notification
type NotificationDevices []NotificationDevice

func (NotificationDevice) TableName() string {
	return "NotificationDevices"
}

func (a *NotificationDevice) BeforeCreate(tx *gorm.DB) (err error) {
	if err := CheckSequence("NotificationDevices", tx); err != nil {
		return err
	}
	return nil
}
