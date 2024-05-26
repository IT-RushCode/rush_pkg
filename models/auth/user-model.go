package auth

import (
	"time"

	"gorm.io/gorm"
)

// Пользователи
type User struct {
	ID           uint           `gorm:"primaryKey;autoincrement"`
	FirstName    string         `gorm:"type:varchar(100)"`
	LastName     string         `gorm:"type:varchar(100)"`
	MiddleName   string         `gorm:"type:varchar(100);default:null"`
	Email        string         `gorm:"type:varchar(100)"`
	PhoneNumber  string         `gorm:"type:varchar(20)"`
	BirthDate    time.Time      `gorm:"type:date"`
	Status       bool           `gorm:"default:true"`
	Avatar       string         `gorm:"type:varchar(255);default:null"`
	UserName     string         `gorm:"type:varchar(100);uniqueIndex"`
	Password     string         `gorm:"type:varchar(255)"`
	IsSuperUser  bool           `gorm:"default:false"`
	LastActivity time.Time      `gorm:"default:null"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Users []User

func (User) TableName() string {
	return "Users"
}
