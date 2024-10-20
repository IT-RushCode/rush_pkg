package models

import "gorm.io/gorm"

// Настройка Юкассы
type YooKassaSetting struct {
	ID        uint   `gorm:"primaryKey"`
	PointID   uint   `gorm:"uniqueIndex;not null"` // Связь с любой моделью (магазин, салон и т.д.)
	StoreID   string `gorm:"not null"`
	SecretKey string `gorm:"not null"`
	IsTest    *bool  `gorm:"deafult:true"`
}

// Настройки Юкассы
type YooKassaSettings []YooKassaSetting

func (YooKassaSetting) TableName() string {
	return "YooKassaSettings"
}

func (a *YooKassaSetting) BeforeCreate(tx *gorm.DB) (err error) {
	if err := CheckSequence("YooKassaSettings", tx); err != nil {
		return err
	}
	return nil
}
