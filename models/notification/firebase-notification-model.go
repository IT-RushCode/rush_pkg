package notification

import "time"

// FirebaseNotification модель
type FirebaseNotification struct {
	ID          uint   `gorm:"primaryKey"`
	DeviceToken string `gorm:"type:varchar(255);not null"`
	Title       string `gorm:"type:varchar(255);not null"`
	Body        string `gorm:"type:text;not null"`
	Data        string `gorm:"type:jsonb"`
	Status      string `gorm:"type:varchar(50);not null;default:'pending'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SentAt      *time.Time
	FailedAt    *time.Time
	Error       string `gorm:"type:text"`
}

type FirebaseNotifications []FirebaseNotification

func (FirebaseNotification) TableName() string {
	return "FirebaseNotifications"
}
