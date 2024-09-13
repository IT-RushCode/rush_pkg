package routes

import (
	"github.com/IT-RushCode/rush_pkg/controllers"
	"github.com/IT-RushCode/rush_pkg/handlers"

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
) {
	yookassa.RUN_YOOKASSA_SETTINGS_ROUTES(api, ctrl)
	yookassa.RUN_PAYMENT_ROUTES(api, h)
}

// РОУТЫ SMS
func RUN_SMS_ROUTES(
	api fiber.Router,
	h *handlers.Handlers,
) {
	sms.RUN_SMS_ROUTES(api, h)
}

// РОУТЫ NOTIFICATION
func RUN_NOTIFICATION_ROUTES(
	api fiber.Router,
	h *handlers.Handlers,
) {
	ntf.RUN_NOTIFICATION_ROUTES(api, h)
}

// РОУТЫ УВЕДОМЛЕНИЙ SMS/EMAIL/PUSH
// func RUN_NOTIFICATION_ROUTES(api fiber.Router, repo *repositories.Repositories) {
// 	go review.RUN_NOTIFICATION(api, repo)
// }

// РОУТЫ GOOGLE КАРТ
// func RUN_MAP_ROUTES(api fiber.Router, repo *repositories.Repositories) {
// 	go review.RUN_MAP(api, repo)
// }