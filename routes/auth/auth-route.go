package routes

import (
	"github.com/IT-RushCode/rush_pkg/config"

	h "github.com/IT-RushCode/rush_pkg/controllers/auth"
	"github.com/IT-RushCode/rush_pkg/repositories"

	"github.com/gofiber/fiber/v2"
)

func RUN_AUTH(api fiber.Router, repo *repositories.Repositories, cfg *config.Config) {
	authController := h.NewAuthController(repo, cfg)

	auth := api.Group("auth")

	auth.Post("/phone-login", authController.PhoneLogin)
	auth.Post("/login", authController.Login)
	auth.Post("/refresh-token", authController.RefreshToken)
	auth.Get("/me", authController.Me)
}
