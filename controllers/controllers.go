package controllers

import (
	"github.com/IT-RushCode/rush_pkg/config"
	ac "github.com/IT-RushCode/rush_pkg/controllers/auth"
	ykc "github.com/IT-RushCode/rush_pkg/controllers/yookassa"
	"github.com/IT-RushCode/rush_pkg/repositories"
)

type Controllers struct {
	// AUTH CONTROLLERS
	AuthController       *ac.AuthController
	UserController       *ac.UserController
	RoleController       *ac.RoleController
	PermissionController *ac.PermissionController

	// YOOKASSA CONTROLLERS
	YookassasettingController *ykc.YookassasettingController
}

// NewControllers - создает новый экземпляр Controllers с инициализированными контроллерами
func NewControllers(repo *repositories.Repositories, cfg *config.Config) *Controllers {
	return &Controllers{
		// AUTH CONTROLLERS
		AuthController:       ac.NewAuthController(repo, cfg),
		UserController:       ac.NewUserController(repo, &cfg.MAIL),
		RoleController:       ac.NewRoleController(repo),
		PermissionController: ac.NewPermissionController(repo),

		// YOOKASSA CONTROLLERS
		YookassasettingController: ykc.NewYooKassaSettingController(repo),
	}
}
