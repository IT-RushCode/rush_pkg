package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`
}

type SoftDelete struct {
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
