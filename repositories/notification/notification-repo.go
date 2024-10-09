package repositories

import (
	"context"
	"time"

	"github.com/IT-RushCode/rush_pkg/models"
	rp "github.com/IT-RushCode/rush_pkg/repositories/base"

	"gorm.io/gorm"
)

// NotificationRepository определяет методы для работы с уведомлениями в базе данных
type NotificationRepository interface {
	rp.BaseRepository
	SaveNotification(ctx context.Context, userID uint, deviceToken, title, message string, notificationType models.NotificationType) error
	ToggleNotificationStatus(ctx context.Context, userID uint, deviceToken string, enable bool) error
	GetNotificationStatus(ctx context.Context, deviceToken string) (bool, error)
	GetEnabledDeviceTokens(ctx context.Context, userID uint) ([]string, error)
	GetNotificationsByUserID(ctx context.Context, userID *uint, filter models.NotificationFilter) ([]models.Notification, error)
}

// notificationRepository - реализация NotificationRepository
type notificationRepository struct {
	rp.BaseRepository
	db *gorm.DB
}

// NewNotificationRepository создает новый экземпляр notificationRepository
func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{
		BaseRepository: rp.NewBaseRepository(db),
		db:             db,
	}
}

// SaveNotification сохраняет уведомление в базу данных
// userID = 0 для общих уведомлений
func (r *notificationRepository) SaveNotification(ctx context.Context, userID uint, deviceToken, title, message string, notificationType models.NotificationType) error {
	now := time.Now()
	notification := models.Notification{
		UserID:  &userID,
		Title:   title,
		Message: message,
		Type:    notificationType, // Указываем тип уведомления
		SentAt:  &now,
	}
	return r.db.WithContext(ctx).Create(&notification).Error
}

// ToggleNotificationStatus обновляет статус уведомлений для указанного устройства
func (r *notificationRepository) ToggleNotificationStatus(ctx context.Context, userID uint, deviceToken string, enable bool) error {
	var notification models.NotificationDevice
	err := r.db.WithContext(ctx).Where("device_token = ?", deviceToken).First(&notification).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Создаем новую запись, если не нашли
			notification = models.NotificationDevice{
				UserID:               &userID,
				DeviceToken:          deviceToken,
				NotificationsEnabled: &enable,
			}
			return r.db.Create(&notification).Error
		}
		return err
	}

	// Если запись найдена, проверяем, нужно ли обновить userID
	if notification.UserID == nil || *notification.UserID == 0 {
		// Если userID не установлен, обновляем его
		if userID != 0 {
			notification.UserID = &userID
		}
	}

	// Обновляем статус уведомлений
	notification.NotificationsEnabled = &enable

	return r.db.Save(&notification).Error
}

// GetNotificationStatus получает текущий статус уведомлений для указанного устройства
func (r *notificationRepository) GetNotificationStatus(ctx context.Context, deviceToken string) (bool, error) {
	var notification models.NotificationDevice
	err := r.db.WithContext(ctx).Where("device_token = ?", deviceToken).First(&notification).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // статус не найден
		}
		return false, err
	}
	return *notification.NotificationsEnabled, nil
}

// GetEnabledDeviceTokens получает все токены устройств с включенными уведомлениями
func (r *notificationRepository) GetEnabledDeviceTokens(ctx context.Context, userID uint) ([]string, error) {
	var tokens []string

	query := r.db.WithContext(ctx).Model(&models.NotificationDevice{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	err := query.Where("notifications_enabled = ?", true).
		Pluck("device_token", &tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// GetNotificationsByUserID возвращает уведомления для пользователя
// Фильтр определяет, какие уведомления возвращать: личные, общие или все
func (r *notificationRepository) GetNotificationsByUserID(ctx context.Context, userID *uint, filter models.NotificationFilter) ([]models.Notification, error) {
	var notifications []models.Notification

	// Создаем базовый запрос
	query := r.db.WithContext(ctx).Model(&models.Notification{})

	switch filter {
	case models.UserNotifications:
		// Возвращаем только личные уведомления
		if userID != nil {
			query = query.Where("user_id = ?", *userID)
		}
	case models.GeneralNotifications:
		// Возвращаем только общие уведомления
		query = query.Where("type = ?", "general")

	case models.AllNotifications:
		// Возвращаем как личные, так и общие уведомления
		query = query.Where("type = ?", "general")
		if userID != nil {
			query = query.Or("user_id = ?", *userID)
		}
	}

	// Выполняем запрос с сортировкой по времени отправки
	err := query.Order("sent_at DESC").Find(&notifications).Error
	if err != nil {
		return nil, err
	}

	return notifications, nil
}
