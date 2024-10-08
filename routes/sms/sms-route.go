package routes

import (
	"github.com/IT-RushCode/rush_pkg/handlers"

	"github.com/gofiber/fiber/v2"
)

func RUN_SMS_ROUTES(api fiber.Router, h *handlers.Handlers) {
	sms := api.Group("sms")

	sms.Post("/send-sms", h.Sms.SendSMS)
	sms.Post("/verify-code", h.Sms.VerifySMSCode)
}
