package auth

import (
	"errors"
	"fmt"
	"time"

	rpBase "github.com/IT-RushCode/rush_pkg/models"
	"github.com/IT-RushCode/rush_pkg/utils"
	"gorm.io/gorm"
)

// Пользователи
type User struct {
	ID             uint      `gorm:"primaryKey;autoincrement"`
	FirstName      string    `gorm:"type:varchar(100)"`
	LastName       string    `gorm:"type:varchar(100)"`
	MiddleName     string    `gorm:"type:varchar(100);default:null"`
	Email          string    `gorm:"type:varchar(100);unique"`
	EmailConfirmed *bool     `gorm:"default:false"`
	PhoneNumber    string    `gorm:"type:varchar(20);unique"`
	PhoneConfirmed *bool     `gorm:"type:varchar(20);unique"`
	BirthDate      time.Time `gorm:"type:date"`
	Status         *bool     `gorm:"default:true"`
	AvatarUrl      string    `gorm:"type:varchar(255);default:null"`
	UserName       string    `gorm:"type:varchar(100);unique"`
	HashPassword   string    `gorm:"type:varchar(255)"`
	IsStaff        *bool     `gorm:"default:false"`
	LastActivity   time.Time `gorm:"default:null"`
	rpBase.BaseModel
}

type Users []User

func (User) TableName() string {
	return "Users"
}

func (user *User) BeforeCreate(db *gorm.DB) (err error) {
	return user.hashPassword()
}

func (user *User) BeforeUpdate(db *gorm.DB) (err error) {
	return user.hashPassword()
}

func (user *User) hashPassword() error {
	hashedPassword, err := utils.HashPassword(user.HashPassword)
	if err != nil {
		return errors.New("ошибка генерации хеша пароля")
	}
	user.HashPassword = hashedPassword
	return nil
}

func (user *User) GetFullName() string {
	return fmt.Sprintf("%s %s %s", user.LastName, user.FirstName, user.MiddleName)
}
