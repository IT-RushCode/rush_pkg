package models

import (
	"time"

	"gorm.io/gorm"
)

// Роль
type Role struct {
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"type:varchar(100);uniqueIndex"`
	Description string         `gorm:"type:varchar(255);default:null"`
	Status      bool           `gorm:"default:true"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// Роли
type Roles []Role

func (Role) TableName() string {
	return "Roles"
}
