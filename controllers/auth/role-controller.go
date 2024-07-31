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

type RoleController struct{ repo *repositories.Repositories }

func NewRoleController(repo *repositories.Repositories) *RoleController {
	return &RoleController{repo: repo}
}

// Создание роли
func (h *RoleController) CreateRole(ctx *fiber.Ctx) error {
	input := &rpAuthDTO.RoleRequestDTO{}
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

	// Привязка привилегий к роли
	if err := h.repo.RolePermission.BindRolePermissions(
		context.Background(),
		data.ID,
		input.Permissions,
	); err != nil {
		return utils.CheckErr(ctx, err)
	}

	role, err := h.repo.Role.FindByIDWithPermissions(context.Background(), data.ID)
	if err != nil {
		return utils.CheckErr(ctx, err)
	}

	return utils.CopyAndRespond(ctx, role, &rpAuthDTO.RoleWithPermissionsDTO{})
}

// Обновление роли
func (h *RoleController) UpdateRole(ctx *fiber.Ctx) error {
	input := &rpAuthDTO.RoleRequestDTO{}
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

	// Привязка привилегий к роли
	if err := h.repo.RolePermission.BindRolePermissions(
		context.Background(),
		data.ID,
		input.Permissions,
	); err != nil {
		return utils.CheckErr(ctx, err)
	}

	role, err := h.repo.Role.FindByIDWithPermissions(context.Background(), data.ID)
	if err != nil {
		return utils.CheckErr(ctx, err)
	}

	return utils.CopyAndRespond(ctx, role, &rpAuthDTO.RoleWithPermissionsDTO{})
}

// Получение роли по ID
func (h *RoleController) FindRoleByID(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data, err := h.repo.Role.FindByIDWithPermissions(context.Background(), id)
	if err != nil {
		return utils.CheckErr(ctx, err)
	}

	return utils.CopyAndRespond(ctx, data, &rpAuthDTO.RoleWithPermissionsDTO{})
}

// Удаление роли
func (h *RoleController) DeleteRole(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data := &rpModels.Role{ID: id}
	if err := h.repo.Role.Delete(context.Background(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	return utils.NoContentResponse(ctx)
}

// Получить все роли с пагинацией или без
func (h *RoleController) GetRoles(ctx *fiber.Ctx) error {
	req, err := utils.GetAllQueries(ctx)
	if err != nil {
		return err
	}

	repoRes := &rpModels.Roles{}
	count, err := h.repo.Role.GetAll(
		context.Background(),
		repoRes,
		req,
	)
	if err != nil {
		return utils.CheckErr(ctx, err)
	}

	resDTO := rpAuthDTO.RolesResponseDTO{}
	if err := copier.Copy(&resDTO, repoRes); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	res := &rpDTO.PaginationDTO{
		List: resDTO,
		Meta: rpDTO.MetaDTO{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
		TotalCount: count,
	}

	return utils.SuccessResponse(ctx, utils.Success, res)
}

// TODO: НУЖНА ДОРАБОТКА

// Получение роли с разрешениями
// func (h *RoleController) GetRoleWithPermissions(ctx *fiber.Ctx) error {
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

// func (h *RoleController) ChangeRolePermission(ctx *fiber.Ctx) error {
// 	return utils.SuccessResponse(ctx, utils.Success, nil)

// }

// func (h *RoleController) GetRolesWithPagination(ctx *fiber.Ctx) error {
// 	return utils.SuccessResponse(ctx, utils.Success, nil)
// }
