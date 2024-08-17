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

type PermissionController struct{ repo *repositories.Repositories }

func NewPermissionController(repo *repositories.Repositories) *PermissionController {
	return &PermissionController{repo: repo}
}

// Создание разрешения
func (h *PermissionController) CreatePermission(ctx *fiber.Ctx) error {
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
func (h *PermissionController) UpdatePermission(ctx *fiber.Ctx) error {
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
func (h *PermissionController) FindPermissionByID(ctx *fiber.Ctx) error {
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
func (h *PermissionController) DeletePermission(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data := &rpModels.Permission{ID: id}
	if err := h.repo.Permission.Delete(context.Background(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	return utils.NoContentResponse(ctx)
}

// Получение всех разрешений с пагинацией или без
func (h *PermissionController) GetPermissions(ctx *fiber.Ctx) error {
	req, err := utils.GetAllQueries(ctx)
	if err != nil {
		return err
	}
	pagination := false

	repoRes := &rpModels.Permissions{}
	count, err := h.repo.Permission.GetAll(
		context.Background(),
		repoRes,
		req,
		pagination,
	)
	if err != nil {
		return utils.CheckErr(ctx, err)
	}
	resDTO := &rpAuthDTO.PermissionsDTO{}
	if err := copier.Copy(resDTO, repoRes); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	if pagination {
		res := &rpDTO.PaginationDTO{
			List: resDTO,
			Meta: rpDTO.MetaDTO{
				Limit:  req.Limit,
				Offset: req.Offset,
			},
			TotalCount: count,
		}

		return utils.SuccessResponse(ctx, utils.Success, res)
	} else {
		return utils.SuccessResponse(ctx, utils.Success, resDTO)
	}

}
