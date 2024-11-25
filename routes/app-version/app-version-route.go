package routes

import (
	"github.com/IT-RushCode/rush_pkg/handlers"
	"github.com/IT-RushCode/rush_pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func RUN_APP_VERSION_ROUTE(app fiber.Router, h *handlers.Handlers, m *middlewares.Middlewares) {
	route := app.Group("app-versions")

	// Получить последнюю версию
	route.Get("/latest", m.Permission.CheckPermission("view:latest_app_vesion"), m.Cache.RouteCache(60), h.AppVersion.GetLatest)

	// Создать новую запись
	route.Post("/", m.Permission.CheckPermission("create:app_vesion"), h.AppVersion.Create)
}
