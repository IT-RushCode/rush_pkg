package routes

import (
	"github.com/IT-RushCode/rush_pkg/controllers"

	"github.com/gofiber/fiber/v2"
)

func RUN_USER(api fiber.Router, ctrl *controllers.Controllers) {
	user := api.Group("users")

	user.Get("/", ctrl.UserController.GetAllUsers)
	user.Get("/:id", ctrl.UserController.FindUserByID)
	user.Post("/", ctrl.UserController.CreateUser)
	user.Put("/:id", ctrl.UserController.UpdateUser)
	user.Delete("/:id", ctrl.UserController.DeleteUser)

	user.Patch("/:id/change-password", ctrl.UserController.ChangeUserPassword)
	user.Patch("/:id/reset-password", ctrl.UserController.ResetUserPassword)

	// user.Patch("/:id/change-status", userController.ChangeUserStatus)
	// user.Patch("/:id/change-roles", userController.ChangeUserRoles)
}
