package routes

import (
	"github.com/IT-RushCode/rush_pkg/handlers"
	"github.com/IT-RushCode/rush_pkg/middlewares"

	appVersion "github.com/IT-RushCode/rush_pkg/routes/app-version"
	chat "github.com/IT-RushCode/rush_pkg/routes/chat"
	ntf "github.com/IT-RushCode/rush_pkg/routes/notification"
	policy "github.com/IT-RushCode/rush_pkg/routes/policy"
	sms "github.com/IT-RushCode/rush_pkg/routes/sms"
	yookassa "github.com/IT-RushCode/rush_pkg/routes/yookassa"

	"github.com/gofiber/fiber/v2"
)

// РОУТЫ ПРОВЕДЕНИЯ ПЛАТЕЖЕЙ ЮКАССЫ
func RUN_YOOKASSA_PAYMENT_ROUTES(
	api fiber.Router,
	h *handlers.Handlers,
	m *middlewares.Middlewares,
) {
	yookassa.RUN_YOOKASSA_SETTINGS_ROUTES(api, h, m)
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
	h *handlers.Handlers,
	m *middlewares.Middlewares,
) {
	ntf.RUN_NOTIFICATION_ROUTES(api, h, m)
}

// РОУТЫ CHAT
func RUN_WEBSOCKET_ROUTES(
	api fiber.Router,
	h *handlers.Handlers,
	m *middlewares.Middlewares,
) {
	chat.RUN_CHAT_API(api, h, m)
	chat.RUN_WEBSOCKET_API(api, h, m)
}

// РОУТЫ NOTIFICATION
func RUN_POLICY_ROUTES(
	api fiber.Router,
	h *handlers.Handlers,
	m *middlewares.Middlewares,
) {
	policy.RUN_POLICY_ROUTES(api, h, m)
}

// РОУТЫ APP VERSION
func RUN_APP_VERSION_ROUTE(
	api fiber.Router,
	h *handlers.Handlers,
	m *middlewares.Middlewares,
) {
	appVersion.RUN_APP_VERSION_ROUTE(api, h, m)
}
