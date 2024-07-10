package routes

import (
	"github.com/IT-RushCode/rush_pkg/config"
	h "github.com/IT-RushCode/rush_pkg/controllers/auth"
	"github.com/IT-RushCode/rush_pkg/repositories"

	"github.com/gofiber/fiber/v2"
)

func RUN_USER(api fiber.Router, repo *repositories.Repositories, cfgMail *config.MailConfig) {
	userController := h.NewUserController(repo, cfgMail)

	user := api.Group("users")

	user.Get("/", userController.GetAllUsers)
	user.Get("/:id", userController.FindUserByID)
	user.Post("/", userController.CreateUser)
	user.Put("/:id", userController.UpdateUser)
	user.Delete("/:id", userController.DeleteUser)

	user.Patch("/:id/change-password", userController.ChangeUserPassword)
	user.Patch("/:id/reset-password", userController.ResetUserPassword)

	// user.Patch("/:id/change-status", userController.ChangeUserStatus)
	// user.Patch("/:id/change-roles", userController.ChangeUserRoles)
}
