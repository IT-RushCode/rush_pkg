package handlers

import (
	"time"

	"github.com/IT-RushCode/rush_pkg/models"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/IT-RushCode/rush_pkg/services"

	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	srv  *services.Services
	repo *repositories.Repositories
}

func NewChatHandler(srv *services.Services, repo *repositories.Repositories) *ChatHandler {
	return &ChatHandler{srv: srv, repo: repo}
}

// POST /api/chat/session
func (h *ChatHandler) CreateChatSession(ctx *fiber.Ctx) error {
	// Предполагаем, что мы получили идентификатор клиента через токен или параметры
	clientID := ctx.Locals("UserID").(uint)

	// Создание новой сессии
	session := models.ChatSession{
		ClientID:  clientID,
		Status:    "active",
		StartedAt: time.Now(),
	}

	// Сохранение сессии в базу данных
	if err := h.repo.Chat.CreateSession(ctx.Context(), &session); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при создании сессии",
		})
	}

	// Возвращаем sessionID клиенту
	return ctx.JSON(fiber.Map{
		"sessionID": session.ID,
	})
}

// GET /api/chat/session
func (h *ChatHandler) GetActiveChatSession(ctx *fiber.Ctx) error {
	clientID := ctx.Locals("UserID").(uint)

	// Получаем активную сессию для клиента
	session, err := h.repo.Chat.GetActiveSession(ctx.Context(), clientID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Активная сессия не найдена",
		})
	}

	// Возвращаем sessionID
	return ctx.JSON(fiber.Map{
		"sessionID": session.ID,
	})
}
