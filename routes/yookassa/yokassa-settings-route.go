package routes

import (
	"github.com/IT-RushCode/rush_pkg/controllers"
	"github.com/gofiber/fiber/v2"
)

func RUN_YOOKASSA_SETTINGS_ROUTES(api fiber.Router, ctrl *controllers.Controllers) {
	ykSetting := api.Group("yookassa-settings")

	ykSetting.Get("/by-point", ctrl.YookassasettingController.FindYooKassaSettingByPointID)
	ykSetting.Get("/:id", ctrl.YookassasettingController.FindYooKassaSettingByID)
	ykSetting.Post("/", ctrl.YookassasettingController.CreateYooKassaSetting)
	ykSetting.Put("/:id", ctrl.YookassasettingController.UpdateYooKassaSetting)
	ykSetting.Delete("/:id", ctrl.YookassasettingController.DeleteYooKassaSetting)
}
