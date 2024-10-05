package models

import (
	"time"

	"gorm.io/gorm"
)

// Notification
type Notification struct {
	ID                   uint             `gorm:"primaryKey;autoIncrement"`
	UserID               *uint            `gorm:"type:varchar(255)"` // ID пользователя
	DeviceToken          *string          `gorm:"type:varchar(255)"` // Токен устройства
	Title                string           `gorm:"type:varchar(255)"` // Заголовок уведомелния
	Message              string           `gorm:"type:text"`         // Текст уведомления
	Type                 NotificationType `gorm:"type:varchar(50)"`  // Новое поле для типа уведомления
	IsGeneral            bool             `gorm:"default:false"`     // Общее или личное уведомление
	NotificationsEnabled *bool            `gorm:"default:true"`      // Статус активности получения уведомления пользователем
	SentAt               time.Time        `gorm:"autoCreateTime"`    // Время отправки
}

// Настройки Notification
type Notifications []Notification

func (Notification) TableName() string {
	return "Notification"
}

func (a *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	// Проверка последовательности id в таблице при автоинкременте
	err = tx.Exec("SELECT setval('\"Notifications_id_seq\"', (SELECT MAX(id) FROM \"Notifications\"));").Error
	if err != nil {
		return err
	}

	return nil
}

// NotificationFilter определяет типы фильтров для выборки уведомлений
type NotificationFilter int

// NotificationType определяет типы уведомлений
type NotificationType string

const (
	// Личные уведомления (по `userID` или `deviceToken`)
	UserNotifications NotificationFilter = iota

	// Общие уведомления
	GeneralNotifications

	// Все уведомления (и личные, и общие)
	AllNotifications

	// Типы уведомлений
	// BirthdayNotification используется для уведомлений, связанных с днями рождения
	BirthdayNotification NotificationType = "birthday"

	// ReminderNotification используется для уведомлений, напоминающих о каком-либо событии
	ReminderNotification NotificationType = "reminder"

	// GeneralNotification используется для общих уведомлений, которые не привязаны к конкретному пользователю
	GeneralNotification NotificationType = "general"

	// PromotionNotification используется для уведомлений о скидках, акциях и специальных предложениях
	PromotionNotification NotificationType = "promotion"
)
