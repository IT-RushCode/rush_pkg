package routes

import (
	"github.com/IT-RushCode/rush_pkg/controllers"
	"github.com/IT-RushCode/rush_pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func RUN_YOOKASSA_SETTINGS_ROUTES(api fiber.Router, ctrl *controllers.Controllers, m *middlewares.Middlewares) {
	ykSetting := api.Group("yookassa-settings")

	ykSetting.Get(
		"/by-point/:pointId",
		m.Permission.CheckPermission("view:yookassa_setting_by_point"),
		m.Cache.RouteCache(60),
		ctrl.YookassaSetting.FindYooKassaSettingByPointID,
	)
	ykSetting.Get("/:id", m.Permission.CheckPermission("view:yookassa_setting_by_id"), m.Cache.RouteCache(60), ctrl.YookassaSetting.FindYooKassaSettingByID)
	ykSetting.Post("/", m.Permission.CheckPermission("create:yookassa_setting"), ctrl.YookassaSetting.CreateYooKassaSetting)
	ykSetting.Put("/:id", m.Permission.CheckPermission("update:yookassa_setting"), ctrl.YookassaSetting.UpdateYooKassaSetting)
	ykSetting.Put("/", m.Permission.CheckPermission("update:yookassa_setting"), ctrl.YookassaSetting.UpdateYooKassaSettingByPointID)
	ykSetting.Delete("/:id", m.Permission.CheckPermission("delete:yookassa_setting"), ctrl.YookassaSetting.DeleteYooKassaSetting)
}
