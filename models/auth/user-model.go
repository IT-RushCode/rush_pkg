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
	ID                      uint      `gorm:"primaryKey;autoincrement"`          // Идентификатор пользователя
	FirstName               string    `gorm:"type:varchar(100)"`                 // Имя
	LastName                string    `gorm:"type:varchar(100)"`                 // Фамилия
	MiddleName              string    `gorm:"type:varchar(100);default:null"`    // Отчество
	Email                   string    `gorm:"type:varchar(100);not null;unique"` // Email
	EmailConfirmed          *bool     `gorm:"default:false"`                     // Подтверждение Email
	PhoneNumber             string    `gorm:"type:varchar(20);unique"`           // Номер телефона
	PhoneConfirmed          *bool     `gorm:"default:false"`                     // Подтверждение номера телефона
	BirthDate               time.Time `gorm:"type:date"`                         // Дата рождения
	Status                  *bool     `gorm:"default:true"`                      // Статус
	AvatarUrl               string    `gorm:"type:varchar(255);default:null"`    // URL аватара
	UserName                string    `gorm:"type:varchar(100);unique"`          // Имя пользователя
	HashPassword            string    `gorm:"type:varchar(255)"`                 // Хеш пароля
	Salt                    string    `gorm:"type:varchar(255)"`                 // Соль для пароля
	ChangePasswordWhenLogin *bool     `gorm:"default:false"`                     // Требование изменения пароля при следующем входе
	IsStaff                 *bool     `gorm:"default:false"`                     // Флаг сотрудника
	LastActivity            time.Time `gorm:"default:null"`                      // Время последней активности

	Roles Roles `gorm:"-"` // Связанные роли, отключение автосоздания таблицы many2many

	rpBase.BaseModel  // Встроенная базовая модель с общими полями
	rpBase.SoftDelete // Мягкое удаление
}

type Users []User

func (User) TableName() string {
	return "Users"
}

func (user *User) BeforeCreate(db *gorm.DB) (err error) {
	salt := utils.GenerateSalt()
	hashedPassword, err := utils.HashPasswordWithSalt(user.HashPassword, salt)
	if err != nil {
		return errors.New("ошибка генерации хеша пароля")
	}
	user.HashPassword = hashedPassword
	user.Salt = salt
	return nil
}

func (user *User) GetFullName() string {
	return fmt.Sprintf("%s %s %s", user.LastName, user.FirstName, user.MiddleName)
}
