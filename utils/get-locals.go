package utils

import (
	"github.com/gofiber/fiber/v2"
)

func GetUserIDFromLocals(ctx *fiber.Ctx) (uint, error) {
	// Получение user_id из локальных данных
	userID, ok := ctx.Locals("UserID").(uint)
	if !ok {
		return 0, ErrorBadRequestResponse(ctx, ErrInvalidData.Error(), nil)
	}

	return userID, nil
}
