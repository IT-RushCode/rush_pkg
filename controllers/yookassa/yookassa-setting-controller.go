package controllers

import (
	dto "github.com/IT-RushCode/rush_pkg/dto/payment"
	ykModel "github.com/IT-RushCode/rush_pkg/models/yookassa"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

type YookassasettingController struct{ repo *repositories.Repositories }

func NewYooKassaSettingController(repo *repositories.Repositories) *YookassasettingController {
	return &YookassasettingController{repo: repo}
}

// Создание YooKassaSetting
func (h *YookassasettingController) CreateYooKassaSetting(ctx *fiber.Ctx) error {
	input := &dto.YooKassaSettingDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	data := &ykModel.YooKassaSetting{}
	if err := copier.Copy(data, input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}
	data.ID = 0

	// TODO: ДОБАВИТЬ ХЕШ ШИФРОВАНИЕ ДЛЯ SECRET KEY
	if err := h.repo.YooKassaSetting.Create(ctx.Context(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	res := &dto.YooKassaSettingDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// Обновление YooKassaSetting
func (h *YookassasettingController) UpdateYooKassaSetting(ctx *fiber.Ctx) error {
	input := &dto.YooKassaSettingDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	data := &ykModel.YooKassaSetting{}
	if err := copier.Copy(data, input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}
	data.ID = id

	// TODO: ДОБАВИТЬ ХЕШ ШИФРОВАНИЕ ДЛЯ SECRET KEY
	if err := h.repo.YooKassaSetting.Update(ctx.Context(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	res := &dto.YooKassaSettingDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// Удаление YooKassaSetting
func (h *YookassasettingController) DeleteYooKassaSetting(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data := &ykModel.YooKassaSetting{ID: id}
	if err := h.repo.YooKassaSetting.Delete(ctx.Context(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	return utils.NoContentResponse(ctx)
}

// Получение разрешения по ID
func (h *YookassasettingController) FindYooKassaSettingByID(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data := &ykModel.YooKassaSetting{}
	if err := h.repo.YooKassaSetting.FindByID(ctx.Context(), id, data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	res := &dto.YooKassaSettingDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// Получение разрешения по ID
func (h *YookassasettingController) FindYooKassaSettingByPointID(ctx *fiber.Ctx) error {

	data := &ykModel.YooKassaSetting{}
	if err := h.repo.YooKassaSetting.Filter(
		ctx.Context(),
		map[string]interface{}{"point_id": uint(ctx.QueryInt("pointId"))},
		data,
	); err != nil {
		return utils.CheckErr(ctx, err)
	}

	res := &dto.YooKassaSettingDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}
