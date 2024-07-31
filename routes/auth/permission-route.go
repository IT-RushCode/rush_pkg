package routes

import (
	"github.com/IT-RushCode/rush_pkg/controllers"
	"github.com/gofiber/fiber/v2"
)

func RUN_PERMISSION(api fiber.Router, ctrl *controllers.Controllers) {
	permission := api.Group("permissions")

	permission.Get("/", ctrl.PermissionController.GetPermissions)
	permission.Get("/:id", ctrl.PermissionController.FindPermissionByID)
	permission.Post("/", ctrl.PermissionController.CreatePermission)
	permission.Put("/:id", ctrl.PermissionController.UpdatePermission)
	permission.Delete("/:id", ctrl.PermissionController.DeletePermission)
}
