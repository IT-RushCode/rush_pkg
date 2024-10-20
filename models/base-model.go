package models

import (
	"fmt"
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

// Проверка последовательности id в таблице при автоинкременте
//
// tableName - принимает название таблицы в базе
//
// Таблица должна иметь поле ID (primary key)
func CheckSequence(tableName string, tx *gorm.DB) error {
	err := tx.Exec(fmt.Sprintf("SELECT setval('\"%s_id_seq\"', (SELECT MAX(id) FROM \"%s\"));", tableName, tableName)).Error
	if err != nil {
		return err
	}
	return nil
}
