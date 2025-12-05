package handlers

import (
	"github.com/IT-RushCode/rush_pkg/config"
	dto "github.com/IT-RushCode/rush_pkg/dto"
	nDto "github.com/IT-RushCode/rush_pkg/dto/notification"
	"github.com/IT-RushCode/rush_pkg/models"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/services"
	"github.com/IT-RushCode/rush_pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
)

type NotificationHandler struct {
	cfg   *config.Config
	srv   *services.Services
	cache *redis.Client
	repo  *repositories.Repositories
}

type notificationStatusDTO struct {
	NotificationStatus bool `json:"notificationStatus"`
}

// NewNotificationHandler создает новый экземпляр NotificationHandler
func NewNotificationHandler(
	cfg *config.Config,
	srv *services.Services,
	cache *redis.Client,
	repo *repositories.Repositories,
) *NotificationHandler {
	return &NotificationHandler{
		cfg:   cfg,
		srv:   srv,
		cache: cache,
		repo:  repo,
	}
}

// SendNotificationsByIdHandler отправляет сохраненное уведомление по идентификатору
// @Summary Отправка сохраненного уведомления
// @Description Отправляет уведомление, которое уже сохранено в системе и подписано пользователями
// @Tags Уведомления
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID уведомления"
// @Success 200 {object} utils.Response "Уведомление отправлено"
// @Failure 400 {object} utils.Response "Некорректные данные"
// @Failure 500 {object} utils.Response "Ошибка сервера"
// @Router /notifications/{id}/send-general [post]
func (h *NotificationHandler) SendNotificationsByIdHandler(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}

	if err := h.srv.Firebase.SendCreatedNotifications(ctx.Context(), uint(id)); err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при отправке уведомлений: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Общие уведомления успешно отправлены", nil)
}

// SendNotificationsHandler обрабатывает отправку уведомления всем пользователям
// @Summary Отправка уведомления всем
// @Description Разрешает отправлять текстовое уведомление одновременно на все устройства
// @Tags Уведомления
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body nDto.SendGeneralNotificationDTO true "Данные уведомления"
// @Success 200 {object} utils.Response "Уведомление отправлено"
// @Failure 400 {object} utils.Response "Некорректные данные"
// @Failure 500 {object} utils.Response "Ошибка сервера"
// @Router /notifications/send-to-all [post]
func (h *NotificationHandler) SendNotificationsHandler(ctx *fiber.Ctx) error {
	var req nDto.SendGeneralNotificationDTO

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := h.srv.Firebase.SendNotifications(ctx.Context(), req.Title, req.Message, models.GeneralNotification); err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при отправке уведомлений: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Общие уведомления успешно отправлены", nil)
}

// SendNotificationToUserHandler отправляет уведомление конкретному пользователю
// @Summary Отправка уведомления конкретному пользователю
// @Description Отправляет уведомление одному пользователю по его ID и типу уведомления
// @Tags Уведомления
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body nDto.SendUserNotificationDTO true "Данные личного уведомления"
// @Success 200 {object} utils.Response "Уведомление отправлено"
// @Failure 400 {object} utils.Response "Некорректные данные"
// @Failure 404 {object} utils.Response "Пользователь не найден"
// @Failure 500 {object} utils.Response "Ошибка сервера"
// @Router /notifications/send-to-user [post]
func (h *NotificationHandler) SendNotificationToUserHandler(ctx *fiber.Ctx) error {
	req := nDto.SendUserNotificationDTO{}
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if req.Type == "general" {
		return utils.ErrorBadRequestResponse(ctx, `поддерживаются только типы - birthday, reminder, promotion`, nil)
	}

	if req.Type == "" {
		return utils.ErrorBadRequestResponse(ctx, "Необходимо указать тип уведомления", nil)
	}

	if err := h.srv.Firebase.SendNotificationToUser(ctx.Context(), req.UserID, req.Title, req.Message, req.Type); err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при отправке личного уведомления: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Личное уведомление успешно отправлено", nil)
}

// ToggleNotificationHandler управляет включением уведомлений на устройстве
// @Summary Изменение статуса уведомлений устройства
// @Description Включает или отключает уведомления для текущего устройства и пользователя
// @Tags Уведомления
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body nDto.ToggleNotificationDTO true "Токен устройства и статус"
// @Success 200 {object} utils.Response "Статус обновлен"
// @Failure 400 {object} utils.Response "Некорректные данные"
// @Failure 500 {object} utils.Response "Ошибка сервера"
// @Router /notifications/toggle [put]
func (h *NotificationHandler) ToggleNotificationHandler(ctx *fiber.Ctx) error {
	userId, err := utils.CheckIsMobile(ctx)
	if err != nil {
		return err
	}

	req := nDto.ToggleNotificationDTO{}
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := h.srv.Firebase.ToggleNotificationStatus(ctx.Context(), userId, req.DeviceToken, req.Enable); err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при обновлении статуса уведомлений: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Статус уведомлений успешно обновлен", nil)
}

// GetToggleNotificationHandler возвращает статус уведомлений по токену устройства
// @Summary Получение статуса уведомлений
// @Description Позволяет узнать, включены ли уведомления для указанного токена
// @Tags Уведомления
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body nDto.GetToggleNotificationDTO true "Токен устройства"
// @Success 200 {object} utils.Response{body=notificationStatusDTO} "Статус уведомлений"
// @Failure 400 {object} utils.Response "Некорректные данные"
// @Failure 500 {object} utils.Response "Ошибка сервера"
// @Router /notifications/get-toggle-status [post]
func (h *NotificationHandler) GetToggleNotificationHandler(ctx *fiber.Ctx) error {
	var req nDto.GetToggleNotificationDTO
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	userId, err := utils.CheckIsMobile(ctx)
	if err != nil {
		return err
	}

	status, err := h.srv.Firebase.GetNotificationStatus(ctx.Context(), userId, req.DeviceToken)
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при получении статуса уведомлений: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Статус уведомлений успешно получен", notificationStatusDTO{NotificationStatus: status})
}

// GetGeneralNotificationsHandler возвращает список общих уведомлений
// @Summary Получение всех общих уведомлений
// @Description Возвращает последние общие уведомления без авторизации
// @Tags Уведомления
// @Produce json
// @Router /notifications/general [get]
// @Success 200 {object} utils.Response "Список уведомлений"
// @Failure 500 {object} utils.Response "Ошибка сервера"
func (h *NotificationHandler) GetGeneralNotificationsHandler(ctx *fiber.Ctx) error {
	notifications, err := h.srv.Firebase.GetNotifications(ctx.Context(), 0, models.GeneralNotifications)
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при получении общих уведомлений: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Общие уведомления успешно получены", mapModelToDTO(notifications))
}

// GetUserNotificationsHandler возвращает личные и общие уведомления для текущего пользователя
// @Summary Получение уведомлений пользователя
// @Description Возвращает уведомления с учетом фильтра: личные, общие или все вместе
// @Tags Уведомления
// @Produce json
// @Security BearerAuth
// @Param filter query int true "Фильтр: 0=личные,1=общие,2=все"
// @Param deviceToken query string false "Токен устройства для фильтрации"
// @Success 200 {object} utils.Response "Уведомления получены"
// @Failure 400 {object} utils.Response "Некорректный фильтр"
// @Failure 500 {object} utils.Response "Ошибка сервера"
// @Router /notifications/user [post]
func (h *NotificationHandler) GetUserNotificationsHandler(ctx *fiber.Ctx) error {
	userId, err := utils.CheckIsMobile(ctx)
	if err != nil {
		return err
	}

	var req nDto.GetNotificationsDTO
	if err := ctx.QueryParser(&req); err != nil {
		return err
	}

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

	notifications, err := h.srv.Firebase.GetNotifications(ctx.Context(), userId, filter)
	if err != nil {
		return utils.ErrorInternalServerErrorResponse(ctx, "Ошибка при получении уведомлений: "+err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Уведомления успешно получены", mapModelToDTO(notifications))
}

// CreateGeneralNotification создает новое общее уведомление
// @Summary Создание общего уведомления
// @Description Создает запись общего уведомления и возвращает DTO
// @Tags Уведомления
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body nDto.NotificationAdminDTO true "Данные уведомления"
// @Success 200 {object} utils.Response{body=nDto.NotificationResponseDTO} "Уведомление создано"
// @Failure 400 {object} utils.Response "Некорректные данные"
// @Failure 500 {object} utils.Response "Ошибка сервера"
// @Router /notifications [post]
func (h *NotificationHandler) CreateGeneralNotification(ctx *fiber.Ctx) error {
	input := &nDto.NotificationAdminDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return err
	}
	if err := utils.ValidateStruct(input); err != nil {
		return err
	}

	data := &models.Notification{Type: models.GeneralNotification}
	if err := copier.Copy(data, input); err != nil {
		return err
	}
	data.ID = 0

	if err := h.repo.Notification.Create(ctx.Context(), data); err != nil {
		return err
	}

	res := &nDto.NotificationResponseDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// UpdateGeneralNotification обновляет существующее общее уведомление
// @Summary Обновление общего уведомления
// @Description Обновляет существующее общее уведомление по ID
// @Tags Уведомления
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID общего уведомления"
// @Param data body nDto.NotificationAdminDTO true "Обновленные данные"
// @Success 200 {object} utils.Response{body=nDto.NotificationResponseDTO} "Уведомление обновлено"
// @Failure 400 {object} utils.Response "Некорректные данные"
// @Failure 404 {object} utils.Response "Не найдено"
// @Failure 500 {object} utils.Response "Ошибка сервера"
// @Router /notifications/{id} [put]
func (h *NotificationHandler) UpdateGeneralNotification(ctx *fiber.Ctx) error {
	input := &nDto.NotificationAdminDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return err
	}
	if err := utils.ValidateStruct(input); err != nil {
		return err
	}

	data := &models.Notification{Type: models.GeneralNotification}
	if err := copier.Copy(data, input); err != nil {
		return err
	}

	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}
	data.ID = id

	if err := h.repo.Notification.Update(ctx.Context(), data); err != nil {
		return err
	}

	res := &nDto.NotificationResponseDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// DeleteNotification удаляет общее уведомление
// @Summary Удаление общего уведомления
// @Description Удаляет уведомление по ID
// @Tags Уведомления
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID уведомления"
// @Success 204 {object} utils.Response "Уведомление удалено"
// @Failure 400 {object} utils.Response "Некорректный ID"
// @Failure 404 {object} utils.Response "Не найдено"
// @Failure 500 {object} utils.Response "Ошибка сервера"
// @Router /notifications/{id} [delete]
func (h *NotificationHandler) DeleteNotification(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data := &models.Notification{ID: id}
	if err := h.repo.Notification.Delete(ctx.Context(), data); err != nil {
		return err
	}

	return utils.NoContentResponse(ctx)
}

// GetGeneralNotificationsAdmin возвращает список общих уведомлений с пагинацией
// @Summary Получение общего списка уведомлений
// @Description Возвращает общий список уведомлений с пагинацией и сортировкой
// @Tags Уведомления
// @Produce json
// @Security BearerAuth
// @Param offset query int false "Смещение"
// @Param limit query int false "Лимит"
// @Param sortBy query string false "Поле сортировки"
// @Param orderBy query string false "Направление сортировки"
// @Success 200 {object} utils.Response{body=dto.PaginationDTO} "Список уведомлений"
// @Failure 400 {object} utils.Response "Неверные параметры"
// @Failure 500 {object} utils.Response "Ошибка сервера"
// @Router /notifications [get]
func (h *NotificationHandler) GetGeneralNotificationsAdmin(ctx *fiber.Ctx) error {
	req, err := utils.GetAllQueries(ctx)
	if err != nil {
		return err
	}

	req.Filters = map[string]string{"type": "general"}
	repoRes := &models.Notifications{}
	count, err := h.repo.Notification.GetAll(
		ctx.Context(),
		repoRes,
		req,
		true,
	)
	if err != nil {
		return err
	}

	resDTO := &[]nDto.NotificationAdminDTO{}
	if err := copier.Copy(resDTO, repoRes); err != nil {
		return err
	}

	res := &dto.PaginationDTO{
		List: resDTO,
		Meta: dto.MetaDTO{
			Limit:      req.Limit,
			Offset:     req.Offset,
			TotalCount: count,
		},
	}

	return utils.SuccessResponse(ctx, "", res)
}

// FindNotificationByID возвращает общее уведомление по ID
// @Summary Получение общего уведомления по ID
// @Description Возвращает детальную информацию по общему уведомлению
// @Tags Уведомления
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID уведомления"
// @Success 200 {object} utils.Response{body=nDto.NotificationAdminDTO} "Уведомление найдено"
// @Failure 400 {object} utils.Response "Некорректный ID"
// @Failure 404 {object} utils.Response "Не найдено"
// @Failure 500 {object} utils.Response "Ошибка сервера"
// @Router /notifications/{id} [get]
func (h *NotificationHandler) FindNotificationByID(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data := &models.Notification{}
	if err := h.repo.Notification.Filter(ctx.Context(), map[string]interface{}{
		"id":   id,
		"type": "general",
	}, data); err != nil {
		return err
	}

	res := &nDto.NotificationAdminDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

func mapModelToDTO(notifications models.Notifications) []nDto.NotificationResponseDTO {
	res := make([]nDto.NotificationResponseDTO, len(notifications))
	for i, n := range notifications {
		res[i] = nDto.NotificationResponseDTO{
			Id:      n.ID,
			Title:   n.Title,
			Message: n.Message,
			Type:    string(n.Type),
			SentAt:  n.SentAt,
		}
	}
	return res
}
