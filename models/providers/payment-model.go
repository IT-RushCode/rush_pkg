package providers

import (
	"time"

	"github.com/IT-RushCode/rush_pkg/models/auth"

	"gorm.io/gorm"
)

// Пластиковые карты
type PayCard struct {
	ID            uint           `gorm:"primaryKey"`
	Number        string         `gorm:"type:varchar(16)"`
	DateMonthYear time.Time      `gorm:"type:date"`
	CVV           string         `gorm:"type:varchar(3)"`
	Cardholder    string         `gorm:"type:varchar(255)"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type PayCards []PayCard

func (PayCard) TableName() string {
	return "PayCards"
}

// Карты пользователей
type UserPayCard struct {
	UserID    uint `gorm:"primaryKey;autoIncrement:false"`
	PayCardID uint `gorm:"primaryKey;autoIncrement:false"`
	IsPrimary bool `gorm:"default:false"`
	User      auth.User
	PayCard   PayCard
}

type UserPayCards []UserPayCard

func (UserPayCard) TableName() string {
	return "UserPayCards"
}
