package providers

import (
	"time"

	rpBase "github.com/IT-RushCode/rush_pkg/models"
	"github.com/IT-RushCode/rush_pkg/models/auth"
)

// Пластиковые карты
type PayCard struct {
	ID            uint       `gorm:"primaryKey"`
	Number        string     `gorm:"type:varchar(16)"`
	DateMonthYear *time.Time `gorm:"type:date"`
	CVV           string     `gorm:"type:varchar(3)"`
	Cardholder    string     `gorm:"type:varchar(255)"`
	rpBase.BaseModel
}

type PayCards []PayCard

func (PayCard) TableName() string {
	return "PayCards"
}

// Карты пользователей
type UserPayCard struct {
	UserID    uint  `gorm:"primaryKey;autoIncrement:false"`
	PayCardID uint  `gorm:"primaryKey;autoIncrement:false"`
	IsPrimary *bool `gorm:"default:false"`
	User      auth.User
	PayCard   PayCard
}

type UserPayCards []UserPayCard

func (UserPayCard) TableName() string {
	return "UserPayCards"
}
