package routes

import (
	c "github.com/IT-RushCode/rush_pkg/controllers/yookassa"
	"github.com/IT-RushCode/rush_pkg/repositories"
	"github.com/gofiber/fiber/v2"
)

func RUN_YOOKASSA_SETTINGS_ROUTES(api fiber.Router, repo *repositories.Repositories) {
	ykSettingController := c.NewYooKassaSettingController(repo)

	ykSetting := api.Group("yookassa-settings")

	ykSetting.Get("/by-point", ykSettingController.FindYooKassaSettingByPointID)
	ykSetting.Get("/:id", ykSettingController.FindYooKassaSettingByID)
	ykSetting.Post("/", ykSettingController.CreateYooKassaSetting)
	ykSetting.Put("/:id", ykSettingController.UpdateYooKassaSetting)
	ykSetting.Delete("/:id", ykSettingController.DeleteYooKassaSetting)
}
