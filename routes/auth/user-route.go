package routes

import (
	h "github.com/IT-RushCode/rush_pkg/handlers/auth"
	"github.com/IT-RushCode/rush_pkg/repositories"

	"github.com/gofiber/fiber/v2"
)

func RUN_USER(api fiber.Router, repo *repositories.Repositories) {
	userHandler := h.NewUserHandler(repo)

	user := api.Group("users")

	user.Get("/", userHandler.GetAllUsers)
	user.Get("/:id", userHandler.FindUserByID)
	user.Post("/", userHandler.CreateUser)
	user.Put("/:id", userHandler.UpdateUser)
	user.Delete("/:id", userHandler.DeleteUser)

	// user.Patch("/:id/change-status", userHandler.ChangeUserStatus)
	// user.Patch("/:id/reset-password", userHandler.ResetUserPassword)
	// user.Patch("/:id/change-password", userHandler.ChangeUserPassword)
	// user.Patch("/:id/change-roles", userHandler.ChangeUserRoles)
}
