package routes

import (
	rpm "github.com/IT-RushCode/rush_pkg/middlewares"

	"github.com/IT-RushCode/rush_pkg/handlers"

	"github.com/gofiber/fiber/v2"
)

func RUN_CHAT_API(
	app fiber.Router,
	h *handlers.Handlers,
	m *rpm.Middlewares,
) {
	// Создание новой сессии
	app.Post("/api/chat/session",
		m.Permission.CheckPermission("write:chat"),
		h.Chat.CreateChatSession)

	// Получение существующей сессии
	app.Get("/api/chat/session",
		m.Permission.CheckPermission("write:chat"),
		h.Chat.GetActiveChatSession)
}
