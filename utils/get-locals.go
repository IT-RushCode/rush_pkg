package utils

import (
	"github.com/gofiber/fiber/v2"
)

func GetUserIDFromLocals(ctx *fiber.Ctx) (int64, error) {
	// Получение user_id из локальных данных
	userID, ok := ctx.Locals("UserID").(int64)
	if !ok {
		return 0, ErrorResponse(ctx, "ошибка сервера", nil)
	}

	return userID, nil
}
