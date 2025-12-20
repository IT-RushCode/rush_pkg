package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time `gorm:"autoCreateTime;comment:Время создания записи" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:Время последнего обновления записи" json:"updatedAt"`
}

type SoftDelete struct {
	DeletedAt gorm.DeletedAt `gorm:"index;comment:Время мягкого удаления"`
}

// Проверка последовательности id в таблице при автоинкременте
//
// tableName - принимает название таблицы в базе
//
// Таблица должна иметь поле ID (primary key)
func CheckSequence(tableName string, db *gorm.DB) error {
	err := db.Exec(fmt.Sprintf(
		"SELECT setval('\"%s_id_seq\"', (SELECT MAX(id) FROM \"%s\"));",
		tableName, tableName,
	)).Error
	if err != nil {
		return err
	}
	return nil
}
