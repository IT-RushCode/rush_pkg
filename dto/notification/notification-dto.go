package dto

import (
	"time"

	"github.com/IT-RushCode/rush_pkg/models"
)

type NotificationAdminDTO struct {
	Id      uint   `json:"id"`
	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
}

type NotificationResponseDTO struct {
	Id      uint       `json:"id"`
	Title   string     `json:"title,omitempty"`
	Message string     `json:"message,omitempty"`
	Type    string     `json:"type,omitempty"`
	SentAt  *time.Time `json:"sentAt,omitempty"`
}

type SendGeneralNotificationDTO struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SendUserNotificationDTO struct {
	UserID  uint                    `json:"userId,omitempty"` // Для личных уведомлений
	Title   string                  `json:"title"`
	Message string                  `json:"message"`
	Type    models.NotificationType `json:"type"` // Тип уведомления
}

type ToggleNotificationDTO struct {
	DeviceToken string `json:"deviceToken" validate:"required"` // Токен устройства
	Enable      bool   `json:"enable"`                          // Статус включения/выключения уведомлений
}

type GetToggleNotificationDTO struct {
	UserID      uint   `json:"userId"`                          // ID пользователя (может быть пустым для анонимных)
	DeviceToken string `json:"deviceToken" validate:"required"` // Токен устройства
}

type GetNotificationsDTO struct {
	UserID      *uint   `json:"userId,omitempty"`      // UserID может быть пустым для анонимных пользователей
	DeviceToken *string `json:"deviceToken,omitempty"` // Токен устройства может быть пустым для фильтрации
	Filter      int     `json:"filter"`                // Тип фильтра: 0 - Личные, 1 - Общие, 2 - Все
}
