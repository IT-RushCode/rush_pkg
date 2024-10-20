package controllers

import (
	"context"

	dto "github.com/IT-RushCode/rush_pkg/dto/chat"
	"github.com/IT-RushCode/rush_pkg/models"
	"github.com/IT-RushCode/rush_pkg/repositories"
	rpu "github.com/IT-RushCode/rush_pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type ChatController struct {
	repo *repositories.Repositories
}

func NewChatController(repo *repositories.Repositories) *ChatController {
	return &ChatController{repo: repo}
}

// ---------------- SESSION ----------------

// Создание новой сессии чата
func (c *ChatController) CreateSession(ctx *fiber.Ctx) error {
	var session models.ChatSession
	if err := ctx.BodyParser(&session); err != nil {
		return err
	}

	if err := c.repo.Chat.CreateSession(ctx.Context(), &session); err != nil {
		return err
	}

	return rpu.CreatedResponse(ctx, rpu.Success, session)
}

// Закрытие сессии чата
func (c *ChatController) CloseSession(ctx *fiber.Ctx) error {
	sessionID := ctx.Params("id")
	if err := c.repo.Chat.CloseSession(ctx.Context(), sessionID); err != nil {
		return err
	}

	return rpu.NoContentResponse(ctx)
}

// Создание сообщения
func (c *ChatController) CreateMessage(ctx context.Context, message *models.ChatMessage) error {
	if err := c.repo.Chat.CreateMessage(ctx, message); err != nil {
		return err
	}
	return nil
}

// Получение всех сообщений для сессии чата
func (c *ChatController) GetHistoryMessages(ctx context.Context, dto *dto.GetChatHistoryDTO, msg *[]models.ChatMessage) error {
	var err error
	*msg, err = c.repo.Chat.GetMessages(ctx, dto)
	if err != nil {
		return err
	}
	return nil
}
