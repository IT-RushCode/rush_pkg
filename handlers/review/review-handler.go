package handlers

import (
	"context"
	"errors"
	"net/http"

	rpDTO "github.com/IT-RushCode/rush_pkg/dto"
	dto "github.com/IT-RushCode/rush_pkg/dto/review"
	models "github.com/IT-RushCode/rush_pkg/models/review"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

type reviewHandler struct{ repo *repositories.Repositories }

func NewReviewHandler(repo *repositories.Repositories) *reviewHandler {
	return &reviewHandler{repo: repo}
}

// Создание пользователя
func (h *reviewHandler) CreateReview(ctx *fiber.Ctx) error {
	input := &dto.ReviewRequestDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	data := &models.Review{}
	if err := copier.Copy(data, &input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}
	data.ID = 0

	if err := h.repo.Review.Create(context.Background(), data); err != nil {
		if errors.Is(err, utils.ErrExists) {
			return utils.SendResponse(ctx, false, err.Error(), nil, http.StatusConflict)
		}
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	res := &dto.ReviewResponseDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// Обновление пользователя
func (h *reviewHandler) UpdateReview(ctx *fiber.Ctx) error {
	input := &dto.ReviewRequestDTO{}
	if err := ctx.BodyParser(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}
	if err := utils.ValidateStruct(input); err != nil {
		return utils.ErrorBadRequestResponse(ctx, err.Error(), nil)
	}

	data := &models.Review{}
	if err := copier.Copy(data, &input); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}
	data.ID = uint(id)

	if err := h.repo.Review.Update(context.Background(), data); err != nil {
		if errors.Is(err, utils.ErrRecordNotFound) {
			return utils.SendResponse(ctx, false, err.Error(), nil, http.StatusNotFound)
		}
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	res := &dto.ReviewResponseDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}

// Удаление пользователя
func (h *reviewHandler) DeleteReview(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return utils.SendResponse(ctx, false, "Invalid ID", nil, http.StatusBadRequest)
	}

	data := &models.Review{ID: uint(id)}
	if err := h.repo.Review.Delete(context.Background(), data); err != nil {
		if errors.Is(err, utils.ErrRecordNotFound) {
			return utils.SendResponse(ctx, false, err.Error(), nil, http.StatusNotFound)
		}
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	return utils.SendResponse(ctx, true, "", nil, http.StatusNoContent)
}

// Получить всех пользователей с пагинацией или без
func (h *reviewHandler) GetAllReviews(ctx *fiber.Ctx) error {
	limit := uint(ctx.QueryInt("limit"))
	offset := uint(ctx.QueryInt("offset"))

	repoRes := &models.Reviews{}
	count, err := h.repo.Review.GetAll(context.Background(), offset, limit, repoRes)
	if err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	resDTO := &dto.ReviewsResponseDTO{}
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

	return utils.SuccessResponse(ctx, "success", res)
}

// Получение разрешения по ID
func (h *reviewHandler) FindReviewByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}

	data := &models.Review{}
	if err := h.repo.Review.FindByID(context.Background(), uint(id), data); err != nil {
		return utils.ErrorResponse(ctx, err.Error(), nil)
	}

	res := &dto.ReviewResponseDTO{}
	return utils.CopyAndRespond(ctx, data, res)
}
