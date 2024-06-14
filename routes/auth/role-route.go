package routes

import (
	h "github.com/IT-RushCode/rush_pkg/controllers/auth"
	"github.com/IT-RushCode/rush_pkg/repositories"

	"github.com/gofiber/fiber/v2"
)

func RUN_ROLE(api fiber.Router, repo *repositories.Repositories) {
	roleController := h.NewRoleController(repo)

	role := api.Group("roles")

	role.Get("/", roleController.GetRoles)
	role.Get("/:id", roleController.FindRoleByID)
	role.Post("/", roleController.CreateRole)
	role.Put("/:id", roleController.UpdateRole)
	role.Delete("/:id", roleController.DeleteRole)
}
