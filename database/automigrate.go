package database

import (
	"log"

	"gorm.io/gorm"

	"github.com/IT-RushCode/rush_pkg/models"
)

// Миграция связанных таблиц Iiko
func MigrateIikoIntegration(conn *gorm.DB) {
	err := conn.Debug().AutoMigrate(
		models.IikoIntegration{},
	)
	if err != nil {
		log.Println(err)
	}
}

// Миграция связанных таблиц Юкассы
func MigrateYooKassaSetting(conn *gorm.DB) {
	err := conn.Debug().AutoMigrate(
		models.YooKassaSetting{},
	)
	if err != nil {
		log.Println(err)
	}
}

// Миграция связанных таблиц уведломлений
func MigrateNotification(conn *gorm.DB) {
	err := conn.Debug().AutoMigrate(
		models.Notification{},
		models.NotificationDevice{},
	)
	if err != nil {
		log.Println(err)
	}
}

// Миграция связанных таблиц чата
func MigrateChat(conn *gorm.DB) {
	err := conn.Debug().AutoMigrate(
		models.ChatSession{},
		models.ChatMessage{},
	)
	if err != nil {
		log.Println(err)
	}
}
