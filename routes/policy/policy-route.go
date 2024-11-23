package routes

import (
	"github.com/IT-RushCode/rush_pkg/handlers"
	"github.com/IT-RushCode/rush_pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

// RUN_POLICY_ROUTES регистрирует маршруты для работы с политиками
func RUN_POLICY_ROUTES(api fiber.Router, h *handlers.Handlers, m *middlewares.Middlewares) {
	policy := api.Group("policy")

	// Получение политики по типу
	policy.Get("/:policyType", m.Permission.CheckPermission("view:policy"), m.Cache.RouteCache(60), h.Policy.GetPolicy)

	// Обновление текста политики
	policy.Patch("/:policyType", m.Permission.CheckPermission("update:policy_text"), h.Policy.UpdatePolicyText)
}
