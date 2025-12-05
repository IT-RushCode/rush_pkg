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

// CreateChatSession запускает новую сессию чата для клиента
// @Summary Создание чат-сессии
// @Description Создает новый чат и возвращает идентификатор сессии
// @Tags Chat
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]uint "sessionID"
// @Failure 500 {object} utils.Response "Ошибка создания сессии"
// @Router /api/chat/session [post]
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

// GetActiveChatSession возвращает активную сессию для пользователя
// @Summary Получение активной сессии чата
// @Description Возвращает идентификатор активной сессии чата по пользователю
// @Tags Chat
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]uint "sessionID"
// @Failure 404 {object} utils.Response "Сессия не найдена"
// @Failure 500 {object} utils.Response "Ошибка сервиса чата"
// @Router /api/chat/session [get]
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
