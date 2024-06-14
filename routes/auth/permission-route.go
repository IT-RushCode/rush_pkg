package routes

import (
	h "github.com/IT-RushCode/rush_pkg/controllers/auth"
	"github.com/IT-RushCode/rush_pkg/repositories"

	"github.com/gofiber/fiber/v2"
)

func RUN_PERMISSION(api fiber.Router, repo *repositories.Repositories) {

	permissionController := h.NewPermissionController(repo)

	permission := api.Group("permissions")

	permission.Get("/", permissionController.GetPermissions)
	permission.Get("/:id", permissionController.FindPermissionByID)
	permission.Post("/", permissionController.CreatePermission)
	permission.Put("/:id", permissionController.UpdatePermission)
	permission.Delete("/:id", permissionController.DeletePermission)
}
