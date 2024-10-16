package services

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/IT-RushCode/rush_pkg/models"
	"github.com/IT-RushCode/rush_pkg/repositories"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

// FirebaseService предоставляет методы для работы с уведомлениями через Firebase
type FirebaseService struct {
	repo     *repositories.Repositories
	firebase *firebase.App
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
		repo:     repo,
		firebase: fbApp,
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

// sendFirebaseNotificationMulti отправляет уведомления нескольким пользователям через Firebase Cloud Messaging с использованием SendAll
func (s *FirebaseService) sendFirebaseNotificationMulti(ctx context.Context, deviceTokens []string, title, message string) error {
	const maxTokens = 500

	client, err := s.firebase.Messaging(ctx)
	if err != nil {
		log.Printf("Ошибка при создании клиента Firebase Messaging: %v", err)
		return err
	}

	// Разбиваем на пачки по 500 токенов
	for i := 0; i < len(deviceTokens); i += maxTokens {
		end := i + maxTokens
		if end > len(deviceTokens) {
			end = len(deviceTokens)
		}

		// Создаем слайс сообщений для текущей пачки токенов
		var messages []*messaging.Message
		for _, token := range deviceTokens[i:end] {
			msg := &messaging.Message{
				Token: token,
				Notification: &messaging.Notification{
					Title: title,
					Body:  message,
				},
			}
			messages = append(messages, msg)
		}

		// Отправляем уведомления с помощью SendAll
		response, err := client.SendAll(ctx, messages)
		if err != nil {
			log.Printf("Ошибка при массовой отправке для пакета с %d по %d: %v", i, end, err)
			return err
		}

		// Логируем успешные и неуспешные отправки
		if response.FailureCount > 0 {
			for idx, resp := range response.Responses {
				if !resp.Success {
					log.Printf("Ошибка при отправке на устройство %s: %v", deviceTokens[i+idx], resp.Error)
				}
			}
		}

		log.Printf("Уведомления успешно отправлены: Успехов: %d, Ошибок: %d", response.SuccessCount, response.FailureCount)
	}

	return nil
}

// sendFirebaseNotificationMulti отправляет уведомления нескольким пользователям через Firebase Cloud Messaging
func (s *FirebaseService) sendFirebaseNotification(ctx context.Context, deviceToken string, title, message string) error {
	client, err := s.firebase.Messaging(ctx)
	if err != nil {
		log.Printf("Ошибка при создании клиента Firebase Messaging: %v", err)
		return err
	}

	response, err := client.Send(ctx, &messaging.Message{
		Token: deviceToken,
		Notification: &messaging.Notification{
			Title: title,
			Body:  message,
		},
	})
	if err != nil {
		log.Printf("Ошибка при отправке уведомления %s: %v", response, err)
		return err
	}

	return nil
}

// ----------- SENDERS -----------

// SendNotifications отправляет общие уведомления которые были ранее созданы всем пользователям с включенными уведомлениями через Firebase
func (s *FirebaseService) SendCreatedNotifications(ctx context.Context, notificationID uint) error {
	// Получаем все устройства с включенными уведомлениями
	deviceTokens, err := s.repo.Notification.GetEnabledDeviceTokens(ctx, 0)
	if err != nil {
		log.Printf("Ошибка при получении токенов устройств: %v", err)
		return err
	}

	notification := &models.Notification{}
	err = s.repo.Notification.FindByID(ctx, notificationID, notification)
	if err != nil {
		log.Printf("Ошибка получения уведомлений: %v", err)
		return err
	}

	// TODO: РЕШИТЬ МУЛЬТИКАСТ ОТПРАВКУ
	// // Отправляем уведомление через Firebase
	// if err := s.sendFirebaseNotificationMulti(ctx, token, notification.Title, notification.Message); err != nil {
	// 	log.Printf("Ошибка при отправке уведомления: %v", err)
	// }

	for _, token := range deviceTokens {
		// Отправляем уведомление через Firebase
		if err := s.sendFirebaseNotification(ctx, token, notification.Title, notification.Message); err != nil {
			log.Printf("Ошибка при отправке уведомления: %v", err)
		}
	}

	err = s.repo.Notification.UpdateField(ctx, notificationID, "sent_at", time.Now(), notification)
	if err != nil {
		log.Printf("Ошибка получения уведомлений: %v", err)
		return err
	}

	log.Println("Общие уведомления отправлены.")
	return nil
}

// SendNotificationToUser отправляет уведомление конкретному пользователю и сохраняет его в базе данных
func (s *FirebaseService) SendNotificationToUser(ctx context.Context, userID uint, title, message string, notificationType models.NotificationType) error {
	if notificationType == "general" {
		// Если уведомление общее, отправляем его всем пользователям
		return s.SendNotifications(ctx, title, message, models.GeneralNotification)
	}

	// Получаем токен устройства с включенными уведомлениями
	deviceTokens, err := s.repo.Notification.GetEnabledDeviceTokens(ctx, userID)
	if err != nil || len(deviceTokens) == 0 {
		log.Printf("Нет доступных токенов для пользователя с userID %d", userID)
		return err
	}

	// TODO: РЕШИТЬ МУЛЬТИКАСТ ОТПРАВКУ
	// // Отправляем уведомление через Firebase для каждого токена
	// if err := s.sendFirebaseNotificationMulti(ctx, deviceTokens, title, message); err != nil {
	// 	log.Printf("Ошибка при отправке уведомления: %v", err)
	// 	return err
	// }

	for _, token := range deviceTokens {
		// Отправляем уведомление через Firebase
		if err := s.sendFirebaseNotification(ctx, token, title, message); err != nil {
			log.Printf("Ошибка при отправке уведомления: %v", err)
		}
	}

	// Сохраняем уведомление в базе данных как личное с указанным типом
	err = s.repo.Notification.SaveNotification(ctx, userID, "", title, message, notificationType)
	if err != nil {
		log.Printf("Ошибка при сохранении уведомления для пользователя %d: %v", userID, err)
		return err
	}

	log.Printf("Личное уведомление отправлено и сохранено для пользователя с userID %d", userID)
	return nil
}

// SendNotifications отправляет общие уведомления всем пользователям с включенными уведомлениями через Firebase
func (s *FirebaseService) SendNotifications(ctx context.Context, title, message string, notificationType models.NotificationType) error {
	// Получаем все устройства с включенными уведомлениями
	deviceTokens, err := s.repo.Notification.GetEnabledDeviceTokens(ctx, 0)
	if err != nil {
		log.Printf("Ошибка при получении токенов устройств: %v", err)
		return err
	}

	// TODO: РЕШИТЬ МУЛЬТИКАСТ ОТПРАВКУ
	// // Отправляем уведомление через Firebase
	// if err := s.sendFirebaseNotificationMulti(ctx, token, title, message); err != nil {
	// 	log.Printf("Ошибка при отправке уведомления: %v", err)
	// }

	for _, token := range deviceTokens {
		// Отправляем уведомление через Firebase
		if err := s.sendFirebaseNotification(ctx, token, title, message); err != nil {
			log.Printf("Ошибка при отправке уведомления: %v", err)
		}
	}

	// Сохраняем уведомление как общее с указанным типом
	err = s.repo.Notification.SaveNotification(ctx, 0, "", title, message, models.GeneralNotification)
	if err != nil {
		log.Printf("Ошибка при сохранении общего уведомления: %v", err)
		return err
	}

	log.Println("Общие уведомления отправлены и сохранены.")
	return nil
}

// ----------- REPO -----------

// ToggleNotificationStatus обновляет статус уведомлений для пользователя
func (s *FirebaseService) ToggleNotificationStatus(ctx context.Context, userID uint, deviceToken string, enable bool) error {
	err := s.repo.Notification.ToggleNotificationStatus(ctx, userID, deviceToken, enable)
	if err != nil {
		log.Printf("Ошибка при обновлении статуса уведомлений: %v", err)
		return err
	}

	log.Printf("Статус уведомлений для устройства %s обновлен на %v", deviceToken, enable)
	return nil
}

// GetNotificationStatus получает текущий статус уведомлений для устройства
func (s *FirebaseService) GetNotificationStatus(ctx context.Context, userID uint, deviceToken string) (bool, error) {
	status, err := s.repo.Notification.GetNotificationStatus(ctx, deviceToken)
	if err != nil {
		log.Printf("Ошибка при получении статуса уведомлений: %v", err)
		return false, err
	}
	return status, nil
}

// GetNotifications возвращает уведомления для пользователя в зависимости от фильтра
// Если userID == "", возвращаем только общие уведомления
func (s *FirebaseService) GetNotifications(ctx context.Context, userID uint, filter models.NotificationFilter) ([]models.Notification, error) {
	// Получаем уведомления с использованием фильтра
	notifications, err := s.repo.Notification.GetNotificationsByUserID(ctx, &userID, filter)
	if err != nil {
		log.Printf("Ошибка при получении уведомлений: %v", err)
		return nil, err
	}

	return notifications, nil
}
