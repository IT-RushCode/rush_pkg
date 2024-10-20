package routes

import (
	"github.com/IT-RushCode/rush_pkg/controllers"
	"github.com/IT-RushCode/rush_pkg/handlers"
	"github.com/IT-RushCode/rush_pkg/middlewares"

	chat "github.com/IT-RushCode/rush_pkg/routes/chat"
	ntf "github.com/IT-RushCode/rush_pkg/routes/notification"
	sms "github.com/IT-RushCode/rush_pkg/routes/sms"
	yookassa "github.com/IT-RushCode/rush_pkg/routes/yookassa"

	"github.com/gofiber/fiber/v2"
)

// РОУТЫ ПРОВЕДЕНИЯ ПЛАТЕЖЕЙ ЮКАССЫ
func RUN_YOOKASSA_PAYMENT_ROUTES(
	api fiber.Router,
	ctrl *controllers.Controllers,
	h *handlers.Handlers,
	m *middlewares.Middlewares,
) {
	yookassa.RUN_YOOKASSA_SETTINGS_ROUTES(api, ctrl, m)
}

// РОУТЫ SMS
func RUN_SMS_ROUTES(
	api fiber.Router,
	h *handlers.Handlers,
	m *middlewares.Middlewares,
) {
	sms.RUN_SMS_ROUTES(api, h)
}

// РОУТЫ NOTIFICATION
func RUN_NOTIFICATION_ROUTES(
	api fiber.Router,
	ctrl *controllers.Controllers,
	h *handlers.Handlers,
	m *middlewares.Middlewares,
) {
	ntf.RUN_NOTIFICATION_ROUTES(api, h, ctrl, m)
}

// РОУТЫ CHAT
func RUN_WEBSOCKET_ROUTES(
	api fiber.Router,
	ctrl *controllers.Controllers,
	h *handlers.Handlers,
	m *middlewares.Middlewares,
) {
	chat.RUN_CHAT_API(api, h, m)
	chat.RUN_WEBSOCKET_API(api, h, m)
}
