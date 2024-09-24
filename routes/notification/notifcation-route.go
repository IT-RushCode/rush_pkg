package routes

import (
	"github.com/IT-RushCode/rush_pkg/handlers"
	"github.com/IT-RushCode/rush_pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RUN_NOTIFICATION_ROUTES(api fiber.Router, h *handlers.Handlers, m *middlewares.Middlewares) {
	sms := api.Group("notifications")

	sms.Post("/send-to-all", m.Permission.CheckPermission("send:notification_to_all"), h.Notification.SendNotificationsHandler)
	sms.Post("/send-to-user", m.Permission.CheckPermission("send:notification_to_user"), h.Notification.SendNotificationToUserHandler)
	sms.Post("/get-toggle-status", h.Notification.GetToggleNotificationHandler)
	sms.Put("/toggle", h.Notification.ToggleNotificationHandler)
}
