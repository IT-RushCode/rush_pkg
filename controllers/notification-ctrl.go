package controllers

import (
	"github.com/IT-RushCode/rush_pkg/dto"
	nDto "github.com/IT-RushCode/rush_pkg/dto/notification"
	"github.com/IT-RushCode/rush_pkg/models"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

type NotificationController struct{ repo *repositories.Repositories }

func NewNotificationController(repo *repositories.Repositories) *NotificationController {
	return &NotificationController{repo: repo}
}

// Создание общего Notification
func (h *NotificationController) CreateGeneralNotification(ctx *fiber.Ctx) error {
	input := &nDto.NotificationAdminDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	data := &models.Notification{Type: models.GeneralNotification}
	if err := copier.Copy(data, input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}
	data.ID = 0

	if err := h.repo.Notification.Create(ctx.Context(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	res := &nDto.NotificationResponseDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// Обновление Notification
func (h *NotificationController) UpdateGeneralNotification(ctx *fiber.Ctx) error {
	input := &nDto.NotificationAdminDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	data := &models.Notification{Type: models.GeneralNotification}
	if err := copier.Copy(data, input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}
	data.ID = id

	if err := h.repo.Notification.Update(ctx.Context(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	res := &nDto.NotificationResponseDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// Удаление Notification
func (h *NotificationController) DeleteNotification(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data := &models.Notification{ID: id}
	if err := h.repo.Notification.Delete(ctx.Context(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	return utils.NoContentResponse(ctx)
}

// Получить все общие Notifications с пагинацией или без
func (h *NotificationController) GetGeneralNotifications(ctx *fiber.Ctx) error {
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
		return utils.CheckErr(ctx, err)
	}

	resDTO := &[]nDto.NotificationAdminDTO{}
	if err := copier.Copy(resDTO, repoRes); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	res := &dto.PaginationDTO{
		List: resDTO,
		Meta: dto.MetaDTO{
			Limit:      req.Limit,
			Offset:     req.Offset,
			TotalCount: count,
		},
	}

	return utils.SuccessResponse(ctx, utils.Success, res)
}

// Получение общий по ID
func (h *NotificationController) FindNotificationByID(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data := &models.Notification{}
	if err := h.repo.Notification.Filter(ctx.Context(), map[string]interface{}{
		"id":   id,
		"type": "general",
	}, data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	res := &nDto.NotificationAdminDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}
