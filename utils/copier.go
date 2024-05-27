package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

func CopyAndRespond(ctx *fiber.Ctx, data interface{}, res interface{}) error {
	if err := copier.Copy(res, data); err != nil {
		return ErrorResponse(ctx, err.Error(), nil)
	}
	return SuccessResponse(ctx, "success", res)
}
