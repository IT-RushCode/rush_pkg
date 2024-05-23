package models

import (
	"time"

	"gorm.io/gorm"
)

// Роли
type Role struct {
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"type:varchar(100);uniqueIndex"`
	Description string         `gorm:"type:varchar(255);default:null"`
	Status      bool           `gorm:"default:true"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Roles []Role

func (Role) TableName() string {
	return "Roles"
}

// Привилегии ролей
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey;autoIncrement:false"`
	PermissionID uint `gorm:"primaryKey;autoIncrement:false"`
	Role         Role
	Permission   Permission
}

type RolePermissions []RolePermission

func (RolePermission) TableName() string {
	return "RolePermissions"
}
