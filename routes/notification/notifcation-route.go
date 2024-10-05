package routes

import (
	"github.com/IT-RushCode/rush_pkg/handlers"
	"github.com/IT-RushCode/rush_pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RUN_NOTIFICATION_ROUTES(api fiber.Router, h *handlers.Handlers, m *middlewares.Middlewares) {
	notifications := api.Group("notifications")

	// Отправка общих уведомлений
	notifications.Post("/send-to-all", m.Permission.CheckPermission("send:notification_to_all"), h.Notification.SendNotificationsHandler)

	// Отправка уведомления конкретному пользователю или общего уведомления
	notifications.Post("/send-to-user", m.Permission.CheckPermission("send:notification_to_user"), h.Notification.SendNotificationToUserHandler)

	// Роут для получения личных и общих уведомлений (с авторизацией)
	notifications.Post("/user", m.Permission.CheckPermission("view:user_notifications"), h.Notification.GetUserNotificationsHandler)

	// Роут для получения только общих уведомлений (без авторизации)
	notifications.Get("/general", h.Notification.GetGeneralNotificationsHandler)

	// Получение статуса уведомлений для пользователя
	notifications.Post("/get-toggle-status", h.Notification.GetToggleNotificationHandler)

	// Управление статусом уведомлений (включение/отключение)
	notifications.Put("/toggle", h.Notification.ToggleNotificationHandler)
}
