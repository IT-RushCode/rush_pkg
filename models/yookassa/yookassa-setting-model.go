package yookassa

import rpBase "github.com/IT-RushCode/rush_pkg/models"

// Настройка Юкассы
type YooKassaSetting struct {
	ID        uint   `gorm:"primaryKey"`
	PointID   uint   `gorm:"uniqueIndex"` // Связь с любой моделью (магазин, салон и т.д.)
	StoreID   string `gorm:"not null"`
	SecretKey string `gorm:"not null"`
	rpBase.BaseModel
}

// Настройки Юкассы
type YooKassaSettings []YooKassaSetting

func (YooKassaSetting) TableName() string {
	return "YooKassaSettings"
}
