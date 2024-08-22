package routes

import (
	"github.com/IT-RushCode/rush_pkg/controllers"
	"github.com/gofiber/fiber/v2"
)

func RUN_YOOKASSA_SETTINGS_ROUTES(api fiber.Router, ctrl *controllers.Controllers) {
	ykSetting := api.Group("yookassa-settings")

	ykSetting.Get("/by-point", ctrl.YookassasettingController.FindYooKassaSettingByPointID).Name("view:yookassa_settings")
	ykSetting.Get("/:id", ctrl.YookassasettingController.FindYooKassaSettingByID).Name("view:yookassa_setting_by_id")
	ykSetting.Post("/", ctrl.YookassasettingController.CreateYooKassaSetting).Name("create:yookassa_setting")
	ykSetting.Put("/:id", ctrl.YookassasettingController.UpdateYooKassaSetting).Name("update:yookassa_setting")
	ykSetting.Delete("/:id", ctrl.YookassasettingController.DeleteYooKassaSetting).Name("delete:yookassa_setting")
}
