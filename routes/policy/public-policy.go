package routes

import (
	"github.com/IT-RushCode/rush_pkg/handlers"
	"github.com/IT-RushCode/rush_pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

// RUN_PUBLIC_ROUTES регистрирует публичные маршруты
func RUN_PUBLIC_ROUTES(app *fiber.App, h *handlers.Handlers, m *middlewares.Middlewares) {
	// Группа маршрутов для публичного доступа
	public := app.Group("/public")

	// Публичный маршрут для политики
	public.Get("/policy/:policyType", m.Cache.RouteCache(60), h.Policy.GetPolicyHTML)
}
