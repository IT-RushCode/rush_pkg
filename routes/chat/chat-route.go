package routes

import (
	rpm "github.com/IT-RushCode/rush_pkg/middlewares"

	"github.com/IT-RushCode/rush_pkg/handlers"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func RUSH_CHAT_API(app fiber.Router, h *handlers.Handlers) {
	// Создание новой сессии
	app.Post("/api/chat/session", h.Chat.CreateChatSession)

	// Получение существующей сессии
	app.Get("/api/chat/session", h.Chat.GetActiveChatSession)
}

// ------------------- WEB SOCKET ------------------->

func RUSH_WEBSOCKET_API(
	app fiber.Router,
	h *handlers.Handlers,
	m *rpm.Middlewares,
) {
	// Маршрут для клиента (мобильное приложение)
	app.Get(
		"/ws/client/chat/:sessionID/:clientID",
		m.Permission.CheckPermission("write:chat"),
		websocket.New(h.WebSocket.WebSocketChat()),
	)

	// Маршрут для поддержки (админка)
	app.Get(
		"/ws/support/chat/:sessionID",
		m.Permission.CheckPermission("manage:chat"),
		websocket.New(h.WebSocket.WebSocketSupport()),
	)
}
