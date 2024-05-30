package handlers

import (
	"context"

	rpDTO "github.com/IT-RushCode/rush_pkg/dto"
	rpAuthDTO "github.com/IT-RushCode/rush_pkg/dto/auth"
	rpModel "github.com/IT-RushCode/rush_pkg/models/auth"
	"github.com/IT-RushCode/rush_pkg/utils"

	"github.com/IT-RushCode/rush_pkg/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

type roleHandler struct{ repo *repositories.Repositories }

func NewRoleHandler(repo *repositories.Repositories) *roleHandler {
	return &roleHandler{repo: repo}
}

// Создание роли
func (h *roleHandler) CreateRole(ctx *fiber.Ctx) error {
	input := &rpAuthDTO.RoleDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	data := &rpModel.Role{}
	if err := copier.Copy(data, &input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}
	data.ID = 0

	if err := h.repo.Role.Create(context.Background(), data); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	res := rpAuthDTO.RoleDTO{}
	return utils.CopyAndRespond(ctx, data, &res)
}

// Обновление роли
func (h *roleHandler) UpdateRole(ctx *fiber.Ctx) error {
	input := &rpAuthDTO.RoleDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	data := &rpModel.Role{}
	if err := copier.Copy(data, &input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}
	data.ID = uint(id)

	if err := h.repo.Role.Update(context.Background(), data); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	res := rpAuthDTO.RoleDTO{}
	return utils.CopyAndRespond(ctx, data, &res)
}

// Получение роли по ID
func (h *roleHandler) FindRoleByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}

	data := &rpModel.Role{}
	if err := h.repo.Role.FindByID(context.Background(), uint(id), data); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	res := rpAuthDTO.RoleDTO{}
	return utils.CopyAndRespond(ctx, data, &res)
}

// Удаление роли
func (h *roleHandler) DeleteRole(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}

	data := &rpModel.Role{ID: uint(id)}
	if err := h.repo.Role.Delete(context.Background(), data); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	return utils.SendResponse(ctx, true, "", nil, fiber.StatusNoContent)
}

// Получить все роли с пагинацией или без
func (h *roleHandler) GetRoles(ctx *fiber.Ctx) error {
	limit, offset := utils.AutoPaginate(ctx)
	repoRes := &rpModel.Roles{}
	count, err := h.repo.Role.GetAll(context.Background(), offset, limit, repoRes)
	if err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
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

	return utils.SuccessResponse(ctx, "success", res)
}

// TODO: НУЖНА ДОРАБОТКА

// Получение роли с разрешениями
// func (h *roleHandler) GetRoleWithPermissions(ctx *fiber.Ctx) error {
// 	id, err := ctx.ParamsInt("id")
// 	if err != nil {
// 		return err
// 	}

// 	data := &rpModel.Role{}
// 	if err := h.repo.Role.FindWithPermissions(context.Background(), uint(id), data); err != nil {
// 		return utils.ErrorResponse(ctx, err.Error(), nil)
// 	}

// 	res := rpAuthDTO.RoleWithPermissionsDTO{}
// 	return utils.CopyAndRespond(ctx, data, &res)
// }

// func (h *roleHandler) ChangeRolePermission(ctx *fiber.Ctx) error {
// 	return utils.SuccessResponse(ctx, "success", nil)

// }

// func (h *roleHandler) GetRolesWithPagination(ctx *fiber.Ctx) error {
// 	return utils.SuccessResponse(ctx, "success", nil)
// }
