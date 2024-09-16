package handlers

import (
	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/IT-RushCode/rush_pkg/services"
	"github.com/IT-RushCode/rush_pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type NotificationHandler struct {
	cfg   *config.Config
	srv   *services.Services
	cache *redis.Client
}

// NewNotificationHandler создает новый экземпляр NotificationHandler
func NewNotificationHandler(
	cfg *config.Config,
	srv *services.Services,
	cache *redis.Client,
) *NotificationHandler {
	return &NotificationHandler{
		cfg:   cfg,
		srv:   srv,
		cache: cache,
	}
}

// SendNotificationsHandler обрабатывает запрос на отправку уведомлений
func (h *NotificationHandler) SendNotificationsHandler(ctx *fiber.Ctx) error {
	var req struct {
		Title   string `json:"title"`
		Message string `json:"message"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, "Ошибка при обработке запроса: "+err.Error(), nil)
	}

	// Вызов сервиса для отправки уведомлений
	err := h.srv.Firebase.SendNotifications(req.Title, req.Message)
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при отправке уведомлений: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Уведомления успешно отправлены", nil)
}

// ToggleNotificationHandler обрабатывает запрос на добавление и управление статусом уведомлений
func (h *NotificationHandler) ToggleNotificationHandler(ctx *fiber.Ctx) error {
	var req struct {
		UserID      string `json:"userId"`      // Поле для авторизованных пользователей, может быть пустым для анонимных
		DeviceToken string `json:"deviceToken"` // Токен устройства
		Enable      bool   `json:"enable"`      // Статус включения/выключения уведомлений
	}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, "Ошибка при обработке запроса: "+err.Error(), nil)
	}

	// Вызов сервиса для обновления статуса уведомлений или добавления токена
	err := h.srv.Firebase.ToggleNotificationStatus(req.UserID, req.DeviceToken, req.Enable)
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при обновлении статуса уведомлений: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Статус уведомлений успешно обновлен", nil)
}

// ToggleNotificationHandler обрабатывает запрос на добавление и управление статусом уведомлений
func (h *NotificationHandler) GetToggleNotificationHandler(ctx *fiber.Ctx) error {
	var req struct {
		UserID      string `json:"userId"`      // Поле для авторизованных пользователей, может быть пустым для анонимных
		DeviceToken string `json:"deviceToken"` // Токен устройства
	}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, "Ошибка при обработке запроса: "+err.Error(), nil)
	}

	// Вызов сервиса для обновления статуса уведомлений или добавления токена
	status, err := h.srv.Firebase.GetToggleNotificationStatus(req.UserID, req.DeviceToken)
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при обновлении статуса уведомлений: "+err.Error(), nil)
	}

	type Status struct {
		NotificationStatus bool `json:"notificationStatus"`
	}

	return utils.SuccessResponse(ctx, utils.Success, Status{NotificationStatus: status})
}
