package routes

import (
	"github.com/IT-RushCode/rush_pkg/handlers"
	"github.com/gofiber/fiber/v2"
)

func RUN_PAYMENT_ROUTES(api fiber.Router, h *handlers.Handlers) {
	payment := api.Group("payment")

	payment.Post("/", h.Payment.CreatePayment)
}
