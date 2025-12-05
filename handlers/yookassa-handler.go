package handlers

import (
	dto "github.com/IT-RushCode/rush_pkg/dto/payment"
	"github.com/IT-RushCode/rush_pkg/models"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

type YooKassaSettingHandler struct {
	repo *repositories.Repositories
}

// NewYooKassaSettingHandler создает новый обработчик для настроек YooKassa
func NewYooKassaSettingHandler(repo *repositories.Repositories) *YooKassaSettingHandler {
	return &YooKassaSettingHandler{repo: repo}
}

// CreateYooKassaSetting создает новую настройку YooKassa
// @Summary Создание настройки YooKassa
// @Description Создает новую запись настроек YooKassa
// @Tags YooKassa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body dto.YooKassaSettingDTO true "Данные конфигурации"
// @Success 200 {object} utils.Response{body=dto.YooKassaSettingDTO} "Конфигурация создана"
// @Failure 400 {object} utils.Response "Некорректные данные"
// @Failure 500 {object} utils.Response "Ошибка сервиса"
// @Router /yookassa-settings [post]
func (h *YooKassaSettingHandler) CreateYooKassaSetting(ctx *fiber.Ctx) error {
	input := &dto.YooKassaSettingDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return err
	}
	if err := utils.ValidateStruct(input); err != nil {
		return err
	}

	data := &models.YooKassaSetting{}
	if err := copier.Copy(data, input); err != nil {
		return err
	}
	data.ID = 0

	if err := h.repo.YooKassaSetting.Create(ctx.Context(), data); err != nil {
		return err
	}

	res := &dto.YooKassaSettingDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// UpdateYooKassaSetting обновляет существующую настройку YooKassa
// @Summary Обновление настройки YooKassa по ID
// @Description Обновляет запись настроек по идентификатору
// @Tags YooKassa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID настройки"
// @Param data body dto.YooKassaSettingDTO true "Обновленные данные"
// @Success 200 {object} utils.Response{body=dto.YooKassaSettingDTO} "Конфигурация обновлена"
// @Failure 400 {object} utils.Response "Некорректные данные"
// @Failure 404 {object} utils.Response "Не найдено"
// @Failure 500 {object} utils.Response "Ошибка сервиса"
// @Router /yookassa-settings/{id} [put]
func (h *YooKassaSettingHandler) UpdateYooKassaSetting(ctx *fiber.Ctx) error {
	input := &dto.YooKassaSettingDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return err
	}
	if err := utils.ValidateStruct(input); err != nil {
		return err
	}

	data := &models.YooKassaSetting{}
	if err := copier.Copy(data, input); err != nil {
		return err
	}

	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}
	data.ID = id

	if err := h.repo.YooKassaSetting.Update(ctx.Context(), data); err != nil {
		return err
	}

	res := &dto.YooKassaSettingDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// UpdateYooKassaSettingByPointID обновляет или создает настройку по pointId
// @Summary Обновление настройки по pointId
// @Description Создает или обновляет настройки YooKassa по идентификатору точки
// @Tags YooKassa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body dto.YooKassaSettingDTO true "Данные конфигурации"
// @Success 200 {object} utils.Response{body=dto.YooKassaSettingDTO} "Конфигурация сохранена"
// @Failure 400 {object} utils.Response "Некорректные данные"
// @Failure 500 {object} utils.Response "Ошибка сервиса"
// @Router /yookassa-settings [put]
func (h *YooKassaSettingHandler) UpdateYooKassaSettingByPointID(ctx *fiber.Ctx) error {
	input := &dto.YooKassaSettingDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return err
	}
	if err := utils.ValidateStruct(input); err != nil {
		return err
	}

	data := &models.YooKassaSetting{}
	if err := copier.Copy(data, input); err != nil {
		return err
	}

	resRepo, err := h.repo.YooKassaSetting.SaveByPointID(ctx.Context(), data)
	if err != nil {
		return err
	}

	res := &dto.YooKassaSettingDTO{}
	return utils.CopyAndRespond(ctx, resRepo, res)
}

// DeleteYooKassaSetting удаляет настройку YooKassa
// @Summary Удаление настройки YooKassa
// @Description Удаляет настройку по идентификатору
// @Tags YooKassa
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID настройки"
// @Success 204 {object} utils.Response "Конфигурация удалена"
// @Failure 400 {object} utils.Response "Некорректный ID"
// @Failure 404 {object} utils.Response "Не найдено"
// @Failure 500 {object} utils.Response "Ошибка сервиса"
// @Router /yookassa-settings/{id} [delete]
func (h *YooKassaSettingHandler) DeleteYooKassaSetting(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data := &models.YooKassaSetting{ID: id}
	if err := h.repo.YooKassaSetting.Delete(ctx.Context(), data); err != nil {
		return err
	}

	return utils.NoContentResponse(ctx)
}

// FindYooKassaSettingByID возвращает настройку по ID
// @Summary Получение настройки по ID
// @Description Возвращает YooKassa настройку по идентификатору
// @Tags YooKassa
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID настройки"
// @Success 200 {object} utils.Response{body=dto.YooKassaSettingDTO} "Конфигурация найдена"
// @Failure 400 {object} utils.Response "Некорректный ID"
// @Failure 404 {object} utils.Response "Не найдено"
// @Failure 500 {object} utils.Response "Ошибка сервиса"
// @Router /yookassa-settings/{id} [get]
func (h *YooKassaSettingHandler) FindYooKassaSettingByID(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data := &models.YooKassaSetting{}
	if err := h.repo.YooKassaSetting.FindByID(ctx.Context(), id, data); err != nil {
		return err
	}

	res := &dto.YooKassaSettingDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// FindYooKassaSettingByPointID возвращает настройку по pointId
// @Summary Получение настройки по pointId
// @Description Возвращает YooKassa настройку по идентификатору точки
// @Tags YooKassa
// @Produce json
// @Security BearerAuth
// @Param pointId path int true "Point ID"
// @Success 200 {object} utils.Response{body=dto.YooKassaSettingDTO} "Конфигурация найдена"
// @Failure 400 {object} utils.Response "Некорректный pointId"
// @Failure 404 {object} utils.Response "Не найдено"
// @Failure 500 {object} utils.Response "Ошибка сервиса"
// @Router /yookassa-settings/by-point/{pointId} [get]
func (h *YooKassaSettingHandler) FindYooKassaSettingByPointID(ctx *fiber.Ctx) error {
	pointID, err := ctx.ParamsInt("pointId")
	if err != nil {
		return err
	}

	data := &models.YooKassaSetting{}
	if err := h.repo.YooKassaSetting.Filter(
		ctx.Context(),
		map[string]interface{}{"point_id": pointID},
		data,
	); err != nil {
		return err
	}

	res := &dto.YooKassaSettingDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}
