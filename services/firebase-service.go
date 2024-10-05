package services

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/IT-RushCode/rush_pkg/models"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"google.golang.org/api/option"
)

// FirebaseService предоставляет методы для работы с уведомлениями через Firebase
type FirebaseService struct {
	Repo     *repositories.Repositories
	Firebase *firebase.App
}

// NewFirebaseService создает новый экземпляр FirebaseService с интеграцией Firebase и локальной базой
func NewFirebaseService(repo *repositories.Repositories, cfg *config.FirebaseConfig) (*FirebaseService, error) {
	credsJSON, err := convertConfigToCredentialsJSON(*cfg)
	if err != nil {
		return nil, err
	}

	opt := option.WithCredentialsJSON(credsJSON)
	fbApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	return &FirebaseService{
		Repo:     repo,
		Firebase: fbApp,
	}, nil
}

// convertConfigToCredentialsJSON формирует JSON для аутентификации Firebase из конфигурации
func convertConfigToCredentialsJSON(cfg config.FirebaseConfig) ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":                        "service_account",
		"project_id":                  cfg.PROJECT_ID,
		"private_key_id":              cfg.PRIVATE_KEY_ID,
		"private_key":                 strings.ReplaceAll(cfg.PRIVATE_KEY, "\\n", "\n"),
		"client_email":                cfg.CLIENT_EMAIL,
		"client_id":                   cfg.CLIENT_ID,
		"auth_uri":                    cfg.AUTH_URI,
		"token_uri":                   cfg.TOKEN_URI,
		"auth_provider_x509_cert_url": cfg.AUTH_PROVIDER_X509_CERT_URL,
		"client_x509_cert_url":        cfg.CLIENT_X509_CERT_URL,
		"universe_domain":             "googleapis.com",
	})
}

// sendFirebaseNotification отправляет уведомление через Firebase Cloud Messaging
func (s *FirebaseService) sendFirebaseNotification(token, title, message string) error {
	ctx := context.Background()
	client, err := s.Firebase.Messaging(ctx)
	if err != nil {
		log.Printf("Ошибка при создании клиента Firebase Messaging: %v", err)
		return err
	}

	// Структура уведомления для отправки
	msg := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  message,
		},
	}

	// Отправка уведомления
	response, err := client.Send(ctx, msg)
	if err != nil {
		log.Printf("Ошибка при отправке уведомления на устройство %s: %v", token, err)
		return err
	}

	log.Printf("Уведомление успешно отправлено на устройство %s: %s", token, response)
	return nil
}

// SendNotificationToUser отправляет уведомление конкретному пользователю и сохраняет его в базе данных
// isGeneral определяет, общее это уведомление или личное
func (s *FirebaseService) SendNotificationToUser(userID uint, title, message string, isGeneral bool, notificationType models.NotificationType) error {
	if isGeneral {
		// Если уведомление общее, отправляем его всем пользователям
		return s.SendNotifications(title, message, notificationType)
	}

	// Проверяем статус уведомлений пользователя
	status, err := s.Repo.Notification.GetNotificationStatus(userID, "")
	if err != nil || !status {
		log.Printf("У пользователя с userID %d нет доступных устройств или уведомления отключены", userID)
		return err
	}

	// Получаем токен устройства с включенными уведомлениями
	deviceTokens, err := s.Repo.Notification.GetEnabledDeviceTokens()
	if err != nil || len(deviceTokens) == 0 {
		log.Printf("Нет доступных токенов для пользователя с userID %d", userID)
		return err
	}

	// Отправляем уведомление через Firebase для каждого токена
	for _, token := range deviceTokens {
		if err := s.sendFirebaseNotification(token, title, message); err != nil {
			log.Printf("Ошибка при отправке уведомления пользователю с userID %d: %v", userID, err)
			return err
		}
	}

	// Сохраняем уведомление в базе данных как личное с указанным типом
	err = s.Repo.Notification.SaveNotification(userID, "", title, message, false, notificationType)
	if err != nil {
		log.Printf("Ошибка при сохранении уведомления для пользователя %d: %v", userID, err)
		return err
	}

	log.Printf("Личное уведомление отправлено и сохранено для пользователя с userID %d", userID)
	return nil
}

// SendNotifications отправляет общие уведомления всем пользователям с включенными уведомлениями через Firebase
func (s *FirebaseService) SendNotifications(title, message string, notificationType models.NotificationType) error {
	// Получаем все устройства с включенными уведомлениями
	tokens, err := s.Repo.Notification.GetEnabledDeviceTokens()
	if err != nil {
		log.Printf("Ошибка при получении токенов устройств: %v", err)
		return err
	}

	for _, token := range tokens {
		// Отправляем уведомление через Firebase
		if err := s.sendFirebaseNotification(token, title, message); err != nil {
			log.Printf("Ошибка при отправке уведомления на устройство с токеном %s: %v", token, err)
		}
	}

	// Сохраняем уведомление как общее с указанным типом
	err = s.Repo.Notification.SaveNotification(0, "", title, message, true, notificationType)
	if err != nil {
		log.Printf("Ошибка при сохранении общего уведомления: %v", err)
		return err
	}

	log.Println("Общие уведомления отправлены и сохранены.")
	return nil
}

// ToggleNotificationStatus обновляет статус уведомлений для пользователя
func (s *FirebaseService) ToggleNotificationStatus(userID uint, deviceToken string, enable bool) error {
	err := s.Repo.Notification.ToggleNotificationStatus(userID, deviceToken, enable)
	if err != nil {
		log.Printf("Ошибка при обновлении статуса уведомлений: %v", err)
		return err
	}

	log.Printf("Статус уведомлений для устройства %s обновлен на %v", deviceToken, enable)
	return nil
}

// GetNotificationStatus получает текущий статус уведомлений для устройства
func (s *FirebaseService) GetNotificationStatus(userID uint, deviceToken string) (bool, error) {
	status, err := s.Repo.Notification.GetNotificationStatus(userID, deviceToken)
	if err != nil {
		log.Printf("Ошибка при получении статуса уведомлений: %v", err)
		return false, err
	}
	return status, nil
}

// GetNotifications возвращает уведомления для пользователя в зависимости от фильтра
// Если userID или deviceToken == "", возвращаем только общие уведомления
func (s *FirebaseService) GetNotifications(userID uint, deviceToken *string, filter models.NotificationFilter) ([]models.Notification, error) {
	// Получаем уведомления с использованием фильтра
	notifications, err := s.Repo.Notification.GetNotificationsByUserID(&userID, deviceToken, filter)
	if err != nil {
		log.Printf("Ошибка при получении уведомлений: %v", err)
		return nil, err
	}

	return notifications, nil
}
