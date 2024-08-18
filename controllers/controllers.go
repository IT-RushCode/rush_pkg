package controllers

import (
	"github.com/IT-RushCode/rush_pkg/config"
	ykc "github.com/IT-RushCode/rush_pkg/controllers/yookassa"
	"github.com/IT-RushCode/rush_pkg/repositories"
)

type Controllers struct {
	// YOOKASSA CONTROLLERS
	YookassasettingController *ykc.YookassasettingController
}

// NewControllers - создает новый экземпляр Controllers с инициализированными контроллерами
func NewControllers(repo *repositories.Repositories, cfg *config.Config) *Controllers {
	return &Controllers{
		// YOOKASSA CONTROLLERS
		YookassasettingController: ykc.NewYooKassaSettingController(repo),
	}
}
