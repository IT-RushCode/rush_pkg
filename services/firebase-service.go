package services

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/IT-RushCode/rush_pkg/config"
)

// FirebaseService предоставляет методы для работы с Firebase
type FirebaseService struct {
	App *firebase.App
}

// NewFirebaseService создает новый экземпляр FirebaseService с данными из конфигурации
func NewFirebaseService(cfg *config.FirebaseConfig) (*FirebaseService, error) {
	credsJSON, err := convertConfigToCredentialsJSON(*cfg)
	if err != nil {
		return nil, err
	}

	opt := option.WithCredentialsJSON(credsJSON)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	return &FirebaseService{App: app}, nil
}

// convertConfigToCredentialsJSON формирует JSON для аутентификации Firebase из конфигурации
func convertConfigToCredentialsJSON(cfg config.FirebaseConfig) ([]byte, error) {
	credentials := map[string]string{
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
	}

	return json.Marshal(credentials)
}

// ToggleNotificationStatus обновляет статус уведомлений для указанного токена устройства
func (s *FirebaseService) ToggleNotificationStatus(userID, deviceToken string, enable bool) error {
	ctx := context.Background()
	client, err := s.App.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	var query firestore.Query

	if userID != "" {
		// Поиск по userID для авторизованных пользователей
		query = client.Collection("users").Where("user_id", "==", userID).Where("token", "==", deviceToken).Limit(1)
	} else {
		// Поиск по токену для анонимных пользователей
		query = client.Collection("users").Where("token", "==", deviceToken).Limit(1)
	}

	iter := query.Documents(ctx)
	doc, err := iter.Next()
	if err == nil && doc.Exists() {
		// Токен существует, обновляем статус уведомлений
		_, err = doc.Ref.Update(ctx, []firestore.Update{
			{Path: "notifications_enabled", Value: enable},
		})
		if err != nil {
			return err
		}
		log.Printf("Статус уведомлений для устройства %s обновлен на %v", deviceToken, enable)
	} else {
		// Токен не найден, добавляем новый документ
		userData := map[string]interface{}{
			"token":                 deviceToken,
			"notifications_enabled": enable,
			"created_at":            time.Now(),
		}
		if userID != "" {
			userData["user_id"] = userID
		}
		_, _, err = client.Collection("users").Add(ctx, userData)
		if err != nil {
			return err
		}
		log.Printf("Новый токен устройства добавлен: %s", deviceToken)
	}

	return nil
}

// UpdateAnonymousToken обновляет или добавляет токен анонимного пользователя в Firestore
func (s *FirebaseService) UpdateAnonymousToken(deviceToken string, enable bool) error {
	ctx := context.Background()
	client, err := s.App.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// Проверяем, существует ли токен в коллекции users
	query := client.Collection("users").Where("token", "==", deviceToken).Limit(1)
	iter := query.Documents(ctx)
	doc, err := iter.Next()
	if err == nil && doc.Exists() {
		// Токен существует, обновляем статус уведомлений
		_, err = doc.Ref.Update(ctx, []firestore.Update{
			{Path: "notifications_enabled", Value: enable},
		})
		if err != nil {
			return err
		}
		log.Printf("Статус уведомлений для анонимного устройства %s обновлен на %v", deviceToken, enable)
	} else {
		// Токен не найден, добавляем новый документ
		_, _, err = client.Collection("anonymous_tokens").Add(ctx, map[string]interface{}{
			"token":                 deviceToken,
			"notifications_enabled": enable,
			"created_at":            time.Now(),
		})
		if err != nil {
			return err
		}
		log.Printf("Новый анонимный токен устройства добавлен: %s", deviceToken)
	}

	return nil
}

// SendNotifications отправляет уведомления всем пользователям с включенными уведомлениями
func (s *FirebaseService) SendNotifications(title, message string) error {
	ctx := context.Background()
	client, err := s.App.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// Поиск всех токенов с включенными уведомлениями
	iter := client.Collection("users").Where("notifications_enabled", "==", true).Documents(ctx)
	messagingClient, err := s.App.Messaging(ctx)
	if err != nil {
		return err
	}

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		token := doc.Data()["token"].(string)
		msg := &messaging.Message{
			Token: token,
			Notification: &messaging.Notification{
				Title: title,
				Body:  message,
			},
		}

		_, err = messagingClient.Send(ctx, msg)
		if err != nil {
			log.Printf("Ошибка при отправке уведомления на токен %s: %v", token, err)
		}
	}

	log.Println("Уведомления отправлены пользователям.")
	return nil
}
