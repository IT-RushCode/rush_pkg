package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

// Функция для маппинга структур
//
// data - принимает исходные данные
//
// res - принимает структуру, в которую нужно поместить данные
func CopyAndRespond(ctx *fiber.Ctx, data interface{}, res interface{}) error {
	if err := copier.Copy(res, data); err != nil {
		return err
	}
	return SuccessResponse(ctx, Success, res)
}
