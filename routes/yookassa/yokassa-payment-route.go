package routes

import (
	h "github.com/IT-RushCode/rush_pkg/handlers"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/gofiber/fiber/v2"
)

func RUN_PAYMENT_ROUTES(api fiber.Router, repo *repositories.Repositories) {
	paymentHandler := h.NewPaymentHandler(repo)

	payment := api.Group("payment")

	payment.Post("/", paymentHandler.CreatePayment)
}
