package routes

import (
	"github.com/IT-RushCode/rush_pkg/handlers"
	"github.com/IT-RushCode/rush_pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func RUN_YOOKASSA_SETTINGS_ROUTES(api fiber.Router, h *handlers.Handlers, m *middlewares.Middlewares) {
	ykSetting := api.Group("yookassa-settings")

	ykSetting.Get(
		"/by-point/:pointId",
		m.Permission.CheckPermission("view:yookassa_setting_by_point"),
		m.Cache.RouteCache(0),
		h.YooKassaSetting.FindYooKassaSettingByPointID,
	)
	ykSetting.Get("/:id", m.Permission.CheckPermission("view:yookassa_setting_by_id"), m.Cache.RouteCache(0), h.YooKassaSetting.FindYooKassaSettingByID)
	ykSetting.Post("/", m.Permission.CheckPermission("create:yookassa_setting"), h.YooKassaSetting.CreateYooKassaSetting)
	ykSetting.Put("/:id", m.Permission.CheckPermission("update:yookassa_setting"), h.YooKassaSetting.UpdateYooKassaSetting)
	ykSetting.Put("/", m.Permission.CheckPermission("update:yookassa_setting"), h.YooKassaSetting.UpdateYooKassaSettingByPointID)
	ykSetting.Delete("/:id", m.Permission.CheckPermission("delete:yookassa_setting"), h.YooKassaSetting.DeleteYooKassaSetting)
}
