package controllers

import (
	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/IT-RushCode/rush_pkg/repositories"
)

type Controllers struct {
	// YOOKASSA CONTROLLERS
	YookassaSetting *YookassasettingController
}

// NewControllers - создает новый экземпляр Controllers с инициализированными контроллерами
func NewControllers(
	repo *repositories.Repositories,
	cfg *config.Config,
) *Controllers {
	return &Controllers{
		// YOOKASSA CONTROLLERS
		YookassaSetting: NewYooKassaSettingController(repo),
	}
}