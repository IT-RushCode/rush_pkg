package routes

import (
	"github.com/IT-RushCode/rush_pkg/config"

	h "github.com/IT-RushCode/rush_pkg/handlers/auth"
	"github.com/IT-RushCode/rush_pkg/repositories"

	"github.com/gofiber/fiber/v2"
)

func RUN_AUTH(api fiber.Router, repo *repositories.Repositories, cfg *config.Config) {
	authHandler := h.NewAuthHandler(repo, cfg)

	auth := api.Group("auth")

	auth.Post("/phone-login", authHandler.PhoneLogin)
	auth.Post("/login", authHandler.Login)
	auth.Post("/me", authHandler.Me)
	auth.Post("/refresh-token", authHandler.RefreshToken)
}
