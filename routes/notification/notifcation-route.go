package routes

import (
	"github.com/IT-RushCode/rush_pkg/handlers"
	"github.com/IT-RushCode/rush_pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RUN_NOTIFICATION_ROUTES(api fiber.Router, h *handlers.Handlers, m *middlewares.Middlewares) {
	notifications := api.Group("notifications")

	// Отправка общих уведомлений
	notifications.Post(
		"/:id/send-general",
		m.Permission.CheckPermission("send:notifications_by_id"),
		h.Notification.SendNotificationsByIdHandler,
	)

	// Отправка общих уведомлений
	notifications.Post(
		"/send-to-all",
		m.Permission.CheckPermission("send:notifications_to_all"),
		h.Notification.SendNotificationsHandler,
	)

	// Отправка уведомления конкретному пользователю или общего уведомления
	notifications.Post(
		"/send-to-user",
		m.Permission.CheckPermission("send:notification_to_user"),
		h.Notification.SendNotificationToUserHandler,
	)

	// Роут для получения личных и общих уведомлений (с авторизацией)
	notifications.Post(
		"/user",
		m.Permission.CheckPermission("view:user_notifications"),
		h.Notification.GetUserNotificationsHandler,
	)

	// Управление статусом уведомлений (включение/отключение)
	notifications.Put(
		"/toggle",
		m.Permission.CheckPermission("change:notification_toggle"),
		h.Notification.ToggleNotificationHandler,
	)

	// Получение статуса уведомлений для пользователя
	notifications.Post(
		"/get-toggle-status",
		m.Permission.CheckPermission("view:notification_toggle_status"),
		h.Notification.GetToggleNotificationHandler,
	)

	// Роут для получения только общих уведомлений (без авторизации)
	notifications.Get(
		"/general",
		m.Cache.RouteCache(0),
		h.Notification.GetGeneralNotificationsHandler,
	)

	// ----------- CRUD ----------->

	notifications.Get(
		"/",
		m.Permission.CheckPermission("view:general_notifications"),
		m.Cache.RouteCache(0),
		h.Notification.GetGeneralNotificationsAdmin,
	)
	notifications.Get(
		"/:id",
		m.Permission.CheckPermission("view:general_notification_by_id"),
		m.Cache.RouteCache(0),
		h.Notification.FindNotificationByID,
	)
	notifications.Post(
		"/",
		m.Permission.CheckPermission("create:general_notification"),
		h.Notification.CreateGeneralNotification,
	)
	notifications.Put(
		"/:id",
		m.Permission.CheckPermission("update:general_notification"),
		h.Notification.UpdateGeneralNotification,
	)
	notifications.Delete(
		"/:id",
		m.Permission.CheckPermission("delete:notification"),
		h.Notification.DeleteNotification,
	)

}
