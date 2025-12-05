package models

import (
	"gorm.io/gorm"
)

type NotificationDevice struct {
	ID                   uint   `gorm:"primaryKey;comment:Первичный ключ устройства"`
	UserID               *uint  `gorm:"index;comment:Ссылка на пользователя"`
	DeviceToken          string `gorm:"unique;comment:Токен устройства для пушей"`
	NotificationsEnabled *bool  `gorm:"default:true;comment:Разрешены ли уведомления"`
	IsAuthenticated      *bool  `gorm:"default:false;comment:Признак авторизованного устройства"`

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
