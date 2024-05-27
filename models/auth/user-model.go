package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/IT-RushCode/rush_pkg/utils"
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

func (user *User) BeforeCreate(db *gorm.DB) (err error) {
	// result, err := utils.HashPassword("admin")
	result, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.New("ошибка генерации хеша пароля")
	}

	user.Password = result

	return nil
}

func (user *User) GetFullName() string {
	return fmt.Sprintf("%s %s %s", user.LastName, user.FirstName, user.MiddleName)
}
