package routes

import (
	"github.com/IT-RushCode/rush_pkg/config"
	h "github.com/IT-RushCode/rush_pkg/handlers"
	"github.com/IT-RushCode/rush_pkg/repositories"

	"github.com/gofiber/fiber/v2"
)

func RUN_SMS_ROUTES(api fiber.Router, cfg *config.Config, repo *repositories.Repositories) {
	smsHandler := h.NewSMSHandler(cfg, repo.Redis)

	sms := api.Group("sms")

	sms.Post("/send-sms", smsHandler.SendSMS)
	sms.Post("/verify-code", smsHandler.VerifySMSCode)
}
