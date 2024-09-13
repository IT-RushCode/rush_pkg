package routes

import (
	"github.com/IT-RushCode/rush_pkg/handlers"

	"github.com/gofiber/fiber/v2"
)

func RUN_NOTIFICATION_ROUTES(api fiber.Router, h *handlers.Handlers) {
	sms := api.Group("notifications")

	sms.Post("/send", h.Notification.SendNotificationsHandler)
	sms.Post("/toggle", h.Notification.ToggleNotificationHandler)
}
