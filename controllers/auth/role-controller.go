package controllers

import (
	"context"

	rpDTO "github.com/IT-RushCode/rush_pkg/dto"
	rpAuthDTO "github.com/IT-RushCode/rush_pkg/dto/auth"
	rpModels "github.com/IT-RushCode/rush_pkg/models/auth"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

type roleController struct{ repo *repositories.Repositories }

func NewRoleController(repo *repositories.Repositories) *roleController {
	return &roleController{repo: repo}
}

// Создание роли
func (h *roleController) CreateRole(ctx *fiber.Ctx) error {
	input := &rpAuthDTO.RoleDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	data := &rpModels.Role{}
	if err := copier.Copy(data, &input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}
	data.ID = 0

	if err := h.repo.Role.Create(context.Background(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	res := rpAuthDTO.RoleDTO{}
	return utils.CopyAndRespond(ctx, data, &res)
}

// Обновление роли
func (h *roleController) UpdateRole(ctx *fiber.Ctx) error {
	input := &rpAuthDTO.RoleDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	data := &rpModels.Role{}
	if err := copier.Copy(data, &input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}
	data.ID = id

	if err := h.repo.Role.Update(context.Background(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	res := rpAuthDTO.RoleDTO{}
	return utils.CopyAndRespond(ctx, data, &res)
}

// Получение роли по ID
func (h *roleController) FindRoleByID(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data := &rpModels.Role{}
	if err := h.repo.Role.FindByID(context.Background(), id, data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	res := &rpAuthDTO.RoleDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// Удаление роли
func (h *roleController) DeleteRole(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data := &rpModels.Role{ID: id}
	if err := h.repo.Role.Delete(context.Background(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	return utils.SendResponse(ctx, true, "", nil, fiber.StatusNoContent)
}

// Получить все роли с пагинацией или без
func (h *roleController) GetRoles(ctx *fiber.Ctx) error {
	limit, offset := utils.AutoPaginate(ctx)

	repoRes := &rpModels.Roles{}
	count, err := h.repo.Role.GetAll(context.Background(), offset, limit, repoRes)
	if err != nil {
		return utils.CheckErr(ctx, err)
	}

	resDTO := rpAuthDTO.RolesDTO{}
	if err := copier.Copy(&resDTO, repoRes); err != nil {
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

// TODO: НУЖНА ДОРАБОТКА

// Получение роли с разрешениями
// func (h *roleController) GetRoleWithPermissions(ctx *fiber.Ctx) error {
// 	id, err := ctx.ParamsInt("id")
// 	if err != nil {
// 		return err
// 	}

// 	data := &rpModels.Role{}
// 	if err := h.repo.Role.FindWithPermissions(context.Background(), id, data); err != nil {
//		return utils.CheckErr(ctx, err)
// 	}

// 	res := rpAuthDTO.RoleWithPermissionsDTO{}
// 	return utils.CopyAndRespond(ctx, data, &res)
// }

// func (h *roleController) ChangeRolePermission(ctx *fiber.Ctx) error {
// 	return utils.SuccessResponse(ctx, utils.Success, nil)

// }

// func (h *roleController) GetRolesWithPagination(ctx *fiber.Ctx) error {
// 	return utils.SuccessResponse(ctx, utils.Success, nil)
// }
