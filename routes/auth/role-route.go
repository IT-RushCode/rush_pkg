package routes

import (
	h "github.com/IT-RushCode/rush_pkg/handlers/auth"
	"github.com/IT-RushCode/rush_pkg/repositories"

	"github.com/gofiber/fiber/v2"
)

func RUN_ROLE(api fiber.Router, repo *repositories.Repositories) {
	roleHandler := h.NewRoleHandler(repo)

	role := api.Group("roles")

	role.Get("/", roleHandler.GetRoles)
	role.Get("/:id", roleHandler.FindRoleByID)
	role.Post("/", roleHandler.CreateRole)
	role.Put("/:id", roleHandler.UpdateRole)
	role.Delete("/:id", roleHandler.DeleteRole)
}
