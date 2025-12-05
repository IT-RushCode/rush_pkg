package models

import "gorm.io/gorm"

// Настройка Юкассы
type YooKassaSetting struct {
	ID        uint   `gorm:"primaryKey;comment:Первичный ключ настройки"`
	PointID   uint   `gorm:"uniqueIndex;not null;comment:Ссылка на точку (магазин, салон и т.п.)"`
	StoreID   string `gorm:"not null;comment:Идентификатор магазина в ЮKassa"`
	SecretKey string `gorm:"not null;comment:Секретный ключ ЮKassa"`
	IsTest    *bool  `gorm:"default:true;comment:Флаг тестовой среды"`
	Status    *bool  `gorm:"default:false;comment:Флаг активности настройки"`
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
