package routes

import (
	"github.com/IT-RushCode/rush_pkg/controllers"
	"github.com/gofiber/fiber/v2"
)

func RUN_YOOKASSA_SETTINGS_ROUTES(api fiber.Router, ctrl *controllers.Controllers) {
	ykSetting := api.Group("yookassa-settings")

	ykSetting.Get("/by-point", ctrl.YookassaSetting.FindYooKassaSettingByPointID).Name("view:yookassa_setting_by_point")

	ykSetting.Get("/:id", ctrl.YookassaSetting.FindYooKassaSettingByID).Name("view:yookassa_setting_by_id")
	ykSetting.Post("/", ctrl.YookassaSetting.CreateYooKassaSetting).Name("create:yookassa_setting")
	ykSetting.Put("/:id", ctrl.YookassaSetting.UpdateYooKassaSetting).Name("update:yookassa_setting")
	ykSetting.Delete("/:id", ctrl.YookassaSetting.DeleteYooKassaSetting).Name("delete:yookassa_setting")
}
