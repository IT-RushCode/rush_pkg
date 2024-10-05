package dto

type SendGeneralNotificationDTO struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SendUserNotificationDTO struct {
	UserID    *uint  `json:"userId,omitempty"` // Для личных уведомлений
	Title     string `json:"title"`
	Message   string `json:"message"`
	IsGeneral bool   `json:"isGeneral"` // Флаг, указывающий, является ли уведомление общим
}

type ToggleNotificationDTO struct {
	UserID      uint   `json:"userId"`      // ID пользователя
	DeviceToken string `json:"deviceToken"` // Токен устройства
	Enable      bool   `json:"enable"`      // Статус включения/выключения уведомлений
}

type GetToggleNotificationDTO struct {
	UserID      uint   `json:"userId"`      // ID пользователя (может быть пустым для анонимных)
	DeviceToken string `json:"deviceToken"` // Токен устройства
}

type GetNotificationsDTO struct {
	UserID      *uint   `json:"userId,omitempty"`      // UserID может быть пустым для анонимных пользователей
	DeviceToken *string `json:"deviceToken,omitempty"` // Токен устройства может быть пустым для фильтрации
	Filter      int     `json:"filter"`                // Тип фильтра: 0 - Личные, 1 - Общие, 2 - Все
}
