package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/IT-RushCode/rush_pkg/config"
	rpDTO "github.com/IT-RushCode/rush_pkg/dto"
	rpAuthDTO "github.com/IT-RushCode/rush_pkg/dto/auth"
	rpModels "github.com/IT-RushCode/rush_pkg/models/auth"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

type userController struct {
	repo *repositories.Repositories
	cfg  *config.MailConfig
}

func NewUserController(
	repo *repositories.Repositories,
	cfg *config.MailConfig,
) *userController {
	return &userController{
		repo: repo,
		cfg:  cfg,
	}
}

// Создание пользователя
func (h *userController) CreateUser(ctx *fiber.Ctx) error {
	input := &rpAuthDTO.UserRequestDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	if input.UserName == "" {
		input.UserName = strings.Split(input.Email, "@")[0]
	}

	data := &rpModels.User{}
	if err := copier.Copy(data, input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}
	data.ID = 0
	changePassword := true
	data.ChangePasswordWhenLogin = &changePassword

	if err := h.repo.User.Create(context.Background(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	// Привязка ролей к пользователю
	if err := h.repo.UserRole.BindUserRoles(context.Background(), data.ID, input.Roles); err != nil {
		return utils.CheckErr(ctx, err)
	}

	user, err := h.repo.User.FindByIDWithRoles(context.Background(), data.ID)
	if err != nil {
		return utils.CheckErr(ctx, err)
	}

	return utils.CopyAndRespond(ctx, user, &rpAuthDTO.UserResponseDTO{})
}

// Обновление пользователя
func (h *userController) UpdateUser(ctx *fiber.Ctx) error {
	input := &rpAuthDTO.UserRequestDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	if input.UserName == "" {
		input.UserName = strings.Split(input.Email, "@")[0]
	}

	data := &rpModels.User{}
	if err := copier.Copy(data, input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}
	data.ID = id

	if err := h.repo.User.Update(context.Background(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	// Привязка ролей к пользователю
	if err := h.repo.UserRole.BindUserRoles(context.Background(), data.ID, input.Roles); err != nil {
		return utils.CheckErr(ctx, err)
	}

	user, err := h.repo.User.FindByIDWithRoles(context.Background(), data.ID)
	if err != nil {
		return utils.CheckErr(ctx, err)
	}

	return utils.CopyAndRespond(ctx, user, &rpAuthDTO.UserResponseDTO{})
}

// Удаление пользователя
func (h *userController) DeleteUser(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data := &rpModels.User{ID: id}
	if err := h.repo.User.Delete(context.Background(), data); err != nil {
		return utils.CheckErr(ctx, err)
	}

	return utils.SendResponse(ctx, true, "", nil, http.StatusNoContent)
}

// Получить всех пользователей с пагинацией или без
func (h *userController) GetAllUsers(ctx *fiber.Ctx) error {
	req, err := utils.GetAllQueries(ctx)
	if err != nil {
		return err
	}

	repoRes := &rpModels.Users{}
	count, err := h.repo.User.GetAll(
		context.Background(),
		repoRes,
		req,
	)
	if err != nil {
		return utils.CheckErr(ctx, err)
	}

	resDTO := &rpAuthDTO.UsersResponseDTO{}
	if err := copier.Copy(resDTO, repoRes); err != nil {
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

// Получение разрешения по ID
func (h *userController) FindUserByID(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	data, err := h.repo.User.FindByIDWithRoles(context.Background(), id)
	if err != nil {
		return utils.CheckErr(ctx, err)
	}

	res := &rpAuthDTO.UserResponseDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// Изменение пароля пользователя
func (h *userController) ChangeUserPassword(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	var input rpAuthDTO.ChangePasswordRequestDTO
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	if err := utils.ValidateStruct(&input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	if err := h.repo.User.ChangePassword(context.Background(), id, input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Пароль успешно изменен", nil)
}

// Сброс пароля пользователя
func (h *userController) ResetUserPassword(ctx *fiber.Ctx) error {
	id, err := utils.GetID(ctx)
	if err != nil {
		return err
	}

	// Генерация нового пароля
	newPassword := utils.GeneratePassword(8, true)

	// Сброс пароля пользователя
	userEmail, err := h.repo.User.ResetPassword(context.Background(), id, newPassword)
	if err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	// Отправка нового пароля пользователю на почту
	if err := h.sendNewPasswordToUser(userEmail, newPassword); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, "Пароль успешно сброшен и отправлен на почту", nil)
}

// sendNewPasswordToUser отправляет новый пароль пользователю на почту
func (h *userController) sendNewPasswordToUser(userEmail string, newPassword string) error {
	fmt.Println("--------", newPassword, "--------")
	subject := fmt.Sprintf("Новый пароль %s", h.cfg.SenderName) // Заголовок почты
	body := fmt.Sprintf("Ваш новый пароль: %s", newPassword)    // Текст почты
	return utils.SendEmail(h.cfg, userEmail, subject, body)
}
