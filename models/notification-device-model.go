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
	// Проверка последовательности id в таблице при автоинкременте
	err = tx.Exec("SELECT setval('\"NotificationDevices_id_seq\"', (SELECT MAX(id) FROM \"NotificationDevices\"));").Error
	if err != nil {
		return err
	}

	return nil
}
