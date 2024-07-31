package routes

import (
	"github.com/IT-RushCode/rush_pkg/controllers"
	"github.com/gofiber/fiber/v2"
)

func RUN_AUTH(api fiber.Router, ctrl *controllers.Controllers) {
	auth := api.Group("auth")

	auth.Post("/email-login", ctrl.AuthController.EmailLogin)
	auth.Post("/phone-login", ctrl.AuthController.PhoneLogin)
	auth.Post("/login", ctrl.AuthController.Login)
	auth.Post("/refresh-token", ctrl.AuthController.RefreshToken)
	auth.Post("/me", ctrl.AuthController.Me)
}
