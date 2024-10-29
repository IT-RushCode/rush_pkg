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
func (h *NotificationHandler) SendNotificationsByIdHandler(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}

	// Вызов сервиса для отправки общих уведомлений с указанным типом
	err = h.srv.Firebase.SendCreatedNotifications(ctx.Context(), uint(id))
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при отправке уведомлений: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Общие уведомления успешно отправлены", nil)
}

// SendNotificationsHandler обрабатывает запрос на отправку общих уведомлений
func (h *NotificationHandler) SendNotificationsHandler(ctx *fiber.Ctx) error {
	var req dto.SendGeneralNotificationDTO

	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, "Ошибка при обработке запроса: "+err.Error(), nil)
	}

	// Вызов сервиса для отправки общих уведомлений с указанным типом
	err := h.srv.Firebase.SendNotifications(ctx.Context(), req.Title, req.Message, models.GeneralNotification)
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при отправке уведомлений: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Общие уведомления успешно отправлены", nil)
}

// SendNotificationToUserHandler обрабатывает запрос на отправку уведомления одному пользователю или общего уведомления
func (h *NotificationHandler) SendNotificationToUserHandler(ctx *fiber.Ctx) error {
	req := dto.SendUserNotificationDTO{}
	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, "Ошибка при обработке запроса: "+err.Error(), nil)
	}

	if req.Type == "general" {
		return utils.ErrorBadRequestResponse(ctx, `поддерживаются только типы - birthday, reminder, promotion`, nil)
	}

	// Проверка типа уведомления
	if req.Type == "" {
		return utils.ErrorBadRequestResponse(ctx, "Необходимо указать тип уведомления", nil)
	}

	// Личное уведомление
	err := h.srv.Firebase.SendNotificationToUser(ctx.Context(), req.UserID, req.Title, req.Message, req.Type)
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при отправке личного уведомления: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Личное уведомление успешно отправлено", nil)
}

// ToggleNotificationHandler обрабатывает запрос на добавление и управление статусом уведомлений
func (h *NotificationHandler) ToggleNotificationHandler(ctx *fiber.Ctx) error {
	// Используем CheckIsMobile для получения userId в зависимости от устройства
	userId, err := utils.CheckIsMobile(ctx)
	if err != nil {
		return err
	}

	req := dto.ToggleNotificationDTO{}
	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorBadRequestResponse(ctx, "Ошибка при обработке запроса: "+err.Error(), nil)
	}

	// Вызов сервиса для обновления статуса уведомлений или добавления токена
	err = h.srv.Firebase.ToggleNotificationStatus(ctx.Context(), userId, req.DeviceToken, req.Enable)
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

	userId, err := utils.CheckIsMobile(ctx)
	if err != nil {
		return err
	}

	// Вызов сервиса для получения статуса уведомлений
	status, err := h.srv.Firebase.GetNotificationStatus(ctx.Context(), userId, req.DeviceToken)
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
	notifications, err := h.srv.Firebase.GetNotifications(ctx.Context(), 0, models.GeneralNotifications)
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при получении общих уведомлений: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Общие уведомления успешно получены", mapModelToDTO(notifications))
}

// GetUserNotificationsHandler обрабатывает запрос на получение личных уведомлений (и общих)
func (h *NotificationHandler) GetUserNotificationsHandler(ctx *fiber.Ctx) error {
	// Используем CheckIsMobile для получения userId в зависимости от устройства
	userId, err := utils.CheckIsMobile(ctx)
	if err != nil {
		return err
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
		req.UserID = &userId
	case 1:
		filter = models.GeneralNotifications
	case 2:
		filter = models.AllNotifications
		req.UserID = &userId
	default:
		return utils.ErrorBadRequestResponse(ctx, "Неверное значение фильтра", nil)
	}

	// Вызов сервиса для получения уведомлений (userId передаем для фильтрации личных)
	notifications, err := h.srv.Firebase.GetNotifications(ctx.Context(), userId, filter)
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при получении уведомлений: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Уведомления успешно получены", mapModelToDTO(notifications))
}

func mapModelToDTO(notifications models.Notifications) []dto.NotificationResponseDTO {
	res := make([]dto.NotificationResponseDTO, len(notifications))
	for i, n := range notifications {
		res[i] = dto.NotificationResponseDTO{
			Id:      n.ID,
			Title:   n.Title,
			Message: n.Message,
			Type:    string(n.Type),
			SentAt:  n.SentAt,
		}
	}
	return res
}
