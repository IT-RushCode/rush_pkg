package controllers

import (
	"context"

	rpDTO "github.com/IT-RushCode/rush_pkg/dto"
	rpAuthDTO "github.com/IT-RushCode/rush_pkg/dto/auth"
	rpModels "github.com/IT-RushCode/rush_pkg/models/auth"
	"github.com/IT-RushCode/rush_pkg/utils"

	"github.com/IT-RushCode/rush_pkg/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

type permissionController struct{ repo *repositories.Repositories }

func NewPermissionController(repo *repositories.Repositories) *permissionController {
	return &permissionController{repo: repo}
}

// Создание разрешения
func (h *permissionController) CreatePermission(ctx *fiber.Ctx) error {
	input := &rpAuthDTO.PermissionDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	data := &rpModels.Permission{}
	if err := copier.Copy(data, input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}
	data.ID = 0

	if err := h.repo.Permission.Create(context.Background(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	res := &rpAuthDTO.PermissionDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// Обновление разрешения
func (h *permissionController) UpdatePermission(ctx *fiber.Ctx) error {
	input := &rpAuthDTO.PermissionDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	data := &rpModels.Permission{}
	if err := copier.Copy(data, input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}
	data.ID = id

	if err := h.repo.Permission.Update(context.Background(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	res := &rpAuthDTO.PermissionDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// Получение разрешения по ID
func (h *permissionController) FindPermissionByID(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data := &rpModels.Permission{}
	if err := h.repo.Permission.FindByID(context.Background(), id, data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	res := &rpAuthDTO.PermissionDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// Удаление разрешения
func (h *permissionController) DeletePermission(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data := &rpModels.Permission{ID: id}
	if err := h.repo.Permission.Delete(context.Background(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	return utils.SendResponse(ctx, true, "", nil, fiber.StatusNoContent)
}

// Получение всех разрешений с пагинацией или без
func (h *permissionController) GetPermissions(ctx *fiber.Ctx) error {
	limit, offset := utils.AutoPaginate(ctx)

	repoRes := &rpModels.Permissions{}
	count, err := h.repo.Permission.GetAll(context.Background(), offset, limit, repoRes)
	if err != nil {
		return utils.CheckErr(ctx, err)
	}
	resDTO := &rpAuthDTO.PermissionsDTO{}
	if err := copier.Copy(resDTO, repoRes); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	res := &rpDTO.PaginationDTO{
		List: resDTO,
		Meta: rpDTO.MetaDTO{
			Limit:  limit,
			Offset: offset,
		},
		TotalCount: count,
	}

	return utils.SuccessResponse(ctx, utils.Success, res)
}
