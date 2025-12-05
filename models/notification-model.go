package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID      uint             `gorm:"primaryKey;autoIncrement;comment:Первичный ключ уведомления"`
	UserID  *uint            `gorm:"default:null;comment:ID пользователя"`
	Title   string           `gorm:"type:varchar(255);comment:Заголовок уведомления"`
	Message string           `gorm:"type:text;comment:Текст уведомления"`
	Type    NotificationType `gorm:"type:varchar(50);comment:Тип уведомления"`
	SentAt  *time.Time       `gorm:"default:null;comment:Время отправки уведомления"`

	BaseModel
}

// Настройки Notification
type Notifications []Notification

func (Notification) TableName() string {
	return "Notifications"
}

func (m *Notification) BeforeCreate(db *gorm.DB) (err error) {
	if err := CheckSequence(m.TableName(), db); err != nil {
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
