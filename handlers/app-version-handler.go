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

func (c *AppVersionHandler) GetLatest(ctx *fiber.Ctx) error {
	version, err := c.repo.AppVersion.GetLatest(ctx.Context())
	if err != nil {
		return utils.CheckErr(ctx, err)
	}

	return utils.SuccessResponse(ctx, utils.Success, version)
}

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
