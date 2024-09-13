package routes

import (
	"github.com/IT-RushCode/rush_pkg/handlers"
	"github.com/IT-RushCode/rush_pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RUN_NOTIFICATION_ROUTES(api fiber.Router, h *handlers.Handlers, m *middlewares.Middlewares) {
	sms := api.Group("notifications")

	sms.Post("/send", m.Permission.CheckPermission("send:notification"), h.Notification.SendNotificationsHandler)
	sms.Post("/toggle", h.Notification.ToggleNotificationHandler)
}
