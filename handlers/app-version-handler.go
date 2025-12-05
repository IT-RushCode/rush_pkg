package handlers

import (
	"github.com/IT-RushCode/rush_pkg/models"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type AppVersionHandler struct {
	repo *repositories.Repositories
}

func NewAppVersionHandler(repo *repositories.Repositories) *AppVersionHandler {
	return &AppVersionHandler{repo: repo}
}

// GetLatest возвращает последнюю версию приложения
// @Summary Получение последней версии
// @Description Возвращает последнюю опубликованную версию приложения
// @Tags AppVersion
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{body=models.AppVersion} "Последняя версия"
// @Failure 500 {object} utils.Response "Ошибка сервиса"
// @Router /app-versions/latest [get]
func (c *AppVersionHandler) GetLatest(ctx *fiber.Ctx) error {
	version, err := c.repo.AppVersion.GetLatest(ctx.Context())
	if err != nil {
		return utils.CheckErr(ctx, err)
	}

	return utils.SuccessResponse(ctx, utils.Success, version)
}

// Create сохраняет новую запись версии
// @Summary Создание версии приложения
// @Description Создает новую запись версии и возвращает сохраненный объект
// @Tags AppVersion
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body models.AppVersion true "Данные версии"
// @Success 200 {object} utils.Response{body=models.AppVersion} "Версия создана"
// @Failure 400 {object} utils.Response "Некорректные данные"
// @Failure 500 {object} utils.Response "Ошибка сервиса"
// @Router /app-versions [post]
func (c *AppVersionHandler) Create(ctx *fiber.Ctx) error {
	var version models.AppVersion
	if err := ctx.BodyParser(&version); err != nil {
		return err
	}

	if err := c.repo.AppVersion.Create(ctx.Context(), &version); err != nil {
		return utils.CheckErr(ctx, err)
	}

	return utils.SuccessResponse(ctx, utils.Success, version)

}
