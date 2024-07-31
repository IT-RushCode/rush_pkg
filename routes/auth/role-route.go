package routes

import (
	"github.com/IT-RushCode/rush_pkg/controllers"

	"github.com/gofiber/fiber/v2"
)

func RUN_ROLE(api fiber.Router, ctrl *controllers.Controllers) {
	role := api.Group("roles")

	role.Get("/", ctrl.RoleController.GetRoles)
	role.Get("/:id", ctrl.RoleController.FindRoleByID)
	role.Post("/", ctrl.RoleController.CreateRole)
	role.Put("/:id", ctrl.RoleController.UpdateRole)
	role.Delete("/:id", ctrl.RoleController.DeleteRole)
}
