package routes

import (
	h "github.com/IT-RushCode/rush_pkg/handlers/auth"
	"github.com/IT-RushCode/rush_pkg/repositories"

	"github.com/gofiber/fiber/v2"
)

func RUN_PERMISSION(api fiber.Router, repo *repositories.Repositories) {

	permissionHandler := h.NewPermissionHandler(repo)

	permission := api.Group("permissions")

	permission.Get("/", permissionHandler.GetPermissions)
	permission.Get("/:id", permissionHandler.FindPermissionByID)
	permission.Post("/", permissionHandler.CreatePermission)
	permission.Put("/:id", permissionHandler.UpdatePermission)
	permission.Delete("/:id", permissionHandler.DeletePermission)
}
