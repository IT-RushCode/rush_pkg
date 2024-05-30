package handlers

import (
	"context"
	"errors"
	"net/http"

	rpDTO "github.com/IT-RushCode/rush_pkg/dto"
	rpAuthDTO "github.com/IT-RushCode/rush_pkg/dto/auth"
	rpModel "github.com/IT-RushCode/rush_pkg/models/auth"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

type userHandler struct{ repo *repositories.Repositories }

func NewUserHandler(repo *repositories.Repositories) *userHandler {
	return &userHandler{repo: repo}
}

// Создание пользователя
func (h *userHandler) CreateUser(ctx *fiber.Ctx) error {
	input := &rpAuthDTO.UserRequestDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	data := &rpModel.User{}
	if err := copier.Copy(data, &input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}
	data.ID = 0

	if err := h.repo.User.Create(context.Background(), data); err != nil {
		if errors.Is(err, utils.ErrExists) {
			return utils.SendResponse(ctx, false, err.Error(), nil, http.StatusConflict)
		}
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	res := rpAuthDTO.UserResponseDTO{}
	return utils.CopyAndRespond(ctx, data, &res)
}

// Обновление пользователя
func (h *userHandler) UpdateUser(ctx *fiber.Ctx) error {
	input := &rpAuthDTO.UserRequestDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	data := &rpModel.User{}
	if err := copier.Copy(data, &input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}
	data.ID = uint(id)

	if err := h.repo.User.Update(context.Background(), data); err != nil {
		if errors.Is(err, utils.ErrRecordNotFound) {
			return utils.SendResponse(ctx, false, err.Error(), nil, http.StatusNotFound)
		}
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	res := rpAuthDTO.UserResponseDTO{}
	return utils.CopyAndRespond(ctx, data, &res)
}

// Удаление пользователя
func (h *userHandler) DeleteUser(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return utils.SendResponse(ctx, false, "Invalid ID", nil, http.StatusBadRequest)
	}

	data := &rpModel.User{ID: uint(id)}
	if err := h.repo.User.Delete(context.Background(), data); err != nil {
		if errors.Is(err, utils.ErrRecordNotFound) {
			return utils.SendResponse(ctx, false, err.Error(), nil, http.StatusNotFound)
		}
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	return utils.SendResponse(ctx, true, "", nil, http.StatusNoContent)
}

// Получить всех пользователей с пагинацией или без
func (h *userHandler) GetAllUsers(ctx *fiber.Ctx) error {
	limit := uint(ctx.QueryInt("limit"))
	offset := uint(ctx.QueryInt("offset"))

	repoRes := &rpModel.Users{}
	count, err := h.repo.User.GetAll(context.Background(), offset, limit, repoRes)
	if err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	resDTO := rpAuthDTO.UsersResponseDTO{}
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

// Получение разрешения по ID
func (h *userHandler) FindUserByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}

	data := &rpModel.User{}
	if err := h.repo.User.FindByID(context.Background(), uint(id), data); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	res := rpAuthDTO.UserResponseDTO{}
	return utils.CopyAndRespond(ctx, data, &res)
}
