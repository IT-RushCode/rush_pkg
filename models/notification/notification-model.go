package notification

import (
	"time"

	rpBase "github.com/IT-RushCode/rush_pkg/models"
)

// Типы уведомлений
const (
	NotificationTypeEmail = "email"
	NotificationTypeSMS   = "sms"
	NotificationTypePush  = "push"
)

// Статусы уведомлений
const (
	NotificationStatusPending = "pending"
	NotificationStatusSent    = "sent"
	NotificationStatusFailed  = "failed"
)

// Notification модель
type Notification struct {
	ID        uint   `gorm:"primaryKey"`
	Type      string `gorm:"type:varchar(50);not null"`
	Recipient string `gorm:"type:varchar(255);not null"`
	Message   string `gorm:"type:text;not null"`
	Status    string `gorm:"type:varchar(50);not null;default:'pending'"`
	SentAt    *time.Time
	FailedAt  *time.Time
	Error     string `gorm:"type:text"`
	rpBase.BaseModel
}

type Notifications []Notification

func (Notification) TableName() string {
	return "Notifications"
}
