package handlers

import (
	"github.com/IT-RushCode/rush_pkg/config"
	dto "github.com/IT-RushCode/rush_pkg/dto/notification"
	"github.com/IT-RushCode/rush_pkg/models"
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

// SendNotificationsHandler обрабатывает запрос на отправку общих уведомлений
func (h *NotificationHandler) SendNotificationsHandler(ctx *fiber.Ctx) error {
	var req dto.SendGeneralNotificationDTO

	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, "Ошибка при обработке запроса: "+err.Error(), nil)
	}

	// Вызов сервиса для отправки общих уведомлений
	err := h.srv.Firebase.SendNotifications(req.Title, req.Message)
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при отправке уведомлений: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Общие уведомления успешно отправлены", nil)
}

// SendNotificationToUserHandler обрабатывает запрос на отправку уведомления одному пользователю или общего уведомления
func (h *NotificationHandler) SendNotificationToUserHandler(ctx *fiber.Ctx) error {
	var req dto.SendUserNotificationDTO

	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, "Ошибка при обработке запроса: "+err.Error(), nil)
	}

	if req.IsGeneral {
		// Общее уведомление
		err := h.srv.Firebase.SendNotifications(req.Title, req.Message)
		if err != nil {
			return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при отправке общего уведомления: "+err.Error(), nil)
		}
		return utils.SuccessResponse(ctx, "Общее уведомление успешно отправлено", nil)
	}

	// Личное уведомление
	err := h.srv.Firebase.SendNotificationToUser(*req.UserID, req.Title, req.Message, req.IsGeneral)
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при отправке личного уведомления: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Личное уведомление успешно отправлено", nil)
}

// ToggleNotificationHandler обрабатывает запрос на добавление и управление статусом уведомлений
func (h *NotificationHandler) ToggleNotificationHandler(ctx *fiber.Ctx) error {
	var req dto.ToggleNotificationDTO

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

// GetToggleNotificationHandler получает текущий статус уведомлений для указанного токена устройства
func (h *NotificationHandler) GetToggleNotificationHandler(ctx *fiber.Ctx) error {
	var req dto.GetToggleNotificationDTO

	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, "Ошибка при обработке запроса: "+err.Error(), nil)
	}

	// Вызов сервиса для получения статуса уведомлений
	status, err := h.srv.Firebase.GetNotificationStatus(req.UserID, req.DeviceToken)
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при получении статуса уведомлений: "+err.Error(), nil)
	}

	type Status struct {
		NotificationStatus bool `json:"notificationStatus"`
	}

	return utils.SuccessResponse(ctx, "Статус уведомлений успешно получен", Status{NotificationStatus: status})
}

// GetGeneralNotificationsHandler обрабатывает запрос на получение только общих уведомлений
func (h *NotificationHandler) GetGeneralNotificationsHandler(ctx *fiber.Ctx) error {
	// Вызов сервиса для получения общих уведомлений
	notifications, err := h.srv.Firebase.GetNotifications(0, nil, models.GeneralNotifications)
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при получении общих уведомлений: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Общие уведомления успешно получены", notifications)
}

// GetUserNotificationsHandler обрабатывает запрос на получение личных уведомлений (и общих)
func (h *NotificationHandler) GetUserNotificationsHandler(ctx *fiber.Ctx) error {
	// Используем CheckIsMobile для получения userId в зависимости от устройства
	userId, err := utils.CheckIsMobile(ctx)
	if err != nil {
		return err // Ошибка будет возвращена из CheckIsMobile в нужном формате
	}

	// Прочитываем фильтр из запроса
	var req dto.GetNotificationsDTO
	if err := ctx.QueryParser(&req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, "Ошибка при обработке запроса: "+err.Error(), nil)
	}

	// Преобразуем числовой фильтр в тип NotificationFilter
	var filter models.NotificationFilter
	switch req.Filter {
	case 0:
		filter = models.UserNotifications
	case 1:
		filter = models.GeneralNotifications
	case 2:
		filter = models.AllNotifications
	default:
		return utils.ErrorBadRequestResponse(ctx, "Неверное значение фильтра", nil)
	}

	// Вызов сервиса для получения уведомлений (userId передаем для фильтрации личных)
	notifications, err := h.srv.Firebase.GetNotifications(userId, req.DeviceToken, filter)
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при получении уведомлений: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Уведомления успешно получены", notifications)
}
