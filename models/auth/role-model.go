package auth

import rpBase "github.com/IT-RushCode/rush_pkg/models"

// Роль
type Role struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(100);uniqueIndex"`
	Description string `gorm:"type:varchar(255);default:null"`
	Status      *bool  `gorm:"default:true"`
	rpBase.BaseModel
}

// Роли
type Roles []Role

func (Role) TableName() string {
	return "Roles"
}
