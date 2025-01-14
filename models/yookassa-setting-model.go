package models

import "gorm.io/gorm"

// Настройка Юкассы
type YooKassaSetting struct {
	ID        uint   `gorm:"primaryKey"`
	PointID   uint   `gorm:"uniqueIndex;not null"` // Связь с любой моделью (магазин, салон и т.д.)
	StoreID   string `gorm:"not null"`
	SecretKey string `gorm:"not null"`
	IsTest    *bool  `gorm:"deafult:true"`
	Status    *bool  `gorm:"default:false"`
}

// Настройки Юкассы
type YooKassaSettings []YooKassaSetting

func (YooKassaSetting) TableName() string {
	return "YooKassaSettings"
}

func (m *YooKassaSetting) BeforeCreate(db *gorm.DB) (err error) {
	if err := CheckSequence(m.TableName(), db); err != nil {
		return err
	}
	return nil
}
