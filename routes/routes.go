package routes

import (
	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/IT-RushCode/rush_pkg/repositories"
	auth "github.com/IT-RushCode/rush_pkg/routes/auth"

	yookassa "github.com/IT-RushCode/rush_pkg/routes/yookassa"

	"github.com/gofiber/fiber/v2"
)

func RUN_AUTH_ROUTES(api fiber.Router, repo *repositories.Repositories, cfg *config.Config) {
	go auth.RUN_AUTH(api, repo, cfg)
	go auth.RUN_USER(api, repo)
	go auth.RUN_ROLE(api, repo)
	go auth.RUN_PERMISSION(api, repo)
}

// РОУТЫ ПРОВЕДЕНИЯ ПЛАТЕЖЕЙ ЮКАССЫ
func RUN_YOOKASSA_PAYMENT_ROUTES(api fiber.Router, repo *repositories.Repositories) {
	go yookassa.RUN_YOOKASSA_SETTINGS_ROUTES(api, repo)
	go yookassa.RUN_PAYMENT_ROUTES(api, repo)
}

// // РОУТЫ ВОЗВРАТОВ ЮКАССЫ
// func RUN_YOOKASSA_REFUND_ROUTES(api fiber.Router, repo *repositories.Repositories) {
// 	go yookassa.RUN_REFUND_KASSA(api, repo)
// }

// // РОУТЫ НАСТРОЕК МАГАЗИНА ЮКАССЫ
// func RUN_YOOKASSA_SETTINGS_ROUTES(api fiber.Router, repo *repositories.Repositories) {
// 	go yookassa.RUN_SETTINGS_KASSA(api, repo)
// }

// РОУТЫ УВЕДОМЛЕНИЙ SMS/EMAIL/PUSH
// func RUN_NOTIFICATION_ROUTES(api fiber.Router, repo *repositories.Repositories) {
// 	go review.RUN_NOTIFICATION(api, repo)
// }

// РОУТЫ ОПЛАТ
// func RUN_PAYMENT_ROUTES(api fiber.Router, repo *repositories.Repositories) {
// 	go review.RUN_PAYMENT(api, repo)
// }

// РОУТЫ GOOGLE КАРТ
// func RUN_MAP_ROUTES(api fiber.Router, repo *repositories.Repositories) {
// 	go review.RUN_MAP(api, repo)
// }

// РОУТЫ ПРОВАЙДЕРОВ
// func RUN_PROVIDER_ROUTES(api fiber.Router, repo *repositories.Repositories) {
// 	go review.RUN_PROVIDER(api, repo)
// }
