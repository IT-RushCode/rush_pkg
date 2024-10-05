package repositories

import (
	"time"

	"github.com/IT-RushCode/rush_pkg/models"
	"gorm.io/gorm"
)

// NotificationRepository определяет методы для работы с уведомлениями в базе данных
type NotificationRepository interface {
	SaveNotification(userID uint, deviceToken, title, message string, isGeneral bool) error
	ToggleNotificationStatus(userID uint, deviceToken string, enable bool) error
	GetNotificationStatus(userID uint, deviceToken string) (bool, error)
	GetEnabledDeviceTokens() ([]string, error)
	GetNotificationsByUserID(userID *uint, deviceToken *string, filter models.NotificationFilter) ([]models.Notification, error)
}

// notificationRepository - реализация NotificationRepository
type notificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository создает новый экземпляр notificationRepository
func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{
		db: db,
	}
}

// SaveNotification сохраняет уведомление в базу данных
func (r *notificationRepository) SaveNotification(userID uint, deviceToken, title, message string, isGeneral bool) error {
	notification := models.Notification{
		UserID:      &userID,
		DeviceToken: &deviceToken,
		Title:       title,
		Message:     message,
		IsGeneral:   isGeneral,
		SentAt:      time.Now(),
	}
	return r.db.Create(&notification).Error
}

// ToggleNotificationStatus обновляет статус уведомлений для указанного устройства
func (r *notificationRepository) ToggleNotificationStatus(userID uint, deviceToken string, enable bool) error {
	var notification models.Notification
	err := r.db.Where("device_token = ? AND user_id = ?", deviceToken, userID).First(&notification).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Создаем новую запись, если не нашли
			notification = models.Notification{
				UserID:               &userID,
				DeviceToken:          &deviceToken,
				NotificationsEnabled: &enable,
				SentAt:               time.Now(),
			}
			return r.db.Create(&notification).Error
		}
		return err
	}

	// Обновляем статус, если запись найдена
	notification.NotificationsEnabled = &enable
	return r.db.Save(&notification).Error
}

// GetNotificationStatus получает текущий статус уведомлений для указанного устройства
func (r *notificationRepository) GetNotificationStatus(userID uint, deviceToken string) (bool, error) {
	var notification models.Notification
	err := r.db.Where("device_token = ? AND user_id = ?", deviceToken, userID).First(&notification).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // статус не найден
		}
		return false, err
	}
	return *notification.NotificationsEnabled, nil
}

// GetEnabledDeviceTokens получает все токены устройств с включенными уведомлениями
func (r *notificationRepository) GetEnabledDeviceTokens() ([]string, error) {
	var tokens []string
	err := r.db.Model(&models.Notification{}).
		Where("notifications_enabled = ?", true).
		Pluck("device_token", &tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// GetNotificationsByUserID возвращает уведомления для пользователя
// Фильтр определяет, какие уведомления возвращать: личные, общие или все
func (r *notificationRepository) GetNotificationsByUserID(userID *uint, deviceToken *string, filter models.NotificationFilter) ([]models.Notification, error) {
	var notifications []models.Notification

	// Создаем базовый запрос
	query := r.db.Model(&models.Notification{})

	switch filter {
	case models.UserNotifications:
		// Возвращаем только личные уведомления
		if userID != nil {
			query = query.Where("user_id = ?", *userID)
		} else if deviceToken != nil {
			query = query.Where("device_token = ?", *deviceToken)
		}

	case models.GeneralNotifications:
		// Возвращаем только общие уведомления
		query = query.Where("is_general = ?", true)

	case models.AllNotifications:
		// Возвращаем как личные, так и общие уведомления
		query = query.Where("is_general = ?", true)
		if userID != nil {
			query = query.Or("user_id = ?", *userID)
		} else if deviceToken != nil {
			query = query.Or("device_token = ?", *deviceToken)
		}
	}

	// Выполняем запрос с сортировкой по времени отправки
	err := query.Order("sent_at DESC").Find(&notifications).Error
	if err != nil {
		return nil, err
	}

	return notifications, nil
}
