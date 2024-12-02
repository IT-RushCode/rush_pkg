package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Проверка, если запрос пришел с мобильного устройства то userID берется из JWT.
// Иначе необходимо передать userID в query.
func CheckIsMobile(ctx *fiber.Ctx) (uint, error) {
	isMob, ok := ctx.Locals("IsMob").(bool)
	if !ok {
		isMob = false
	}

	if isMob {
		userId, err := GetUserIDFromLocals(ctx)
		if err != nil {
			return 0, fiber.NewError(fiber.StatusBadRequest, "Ошибка получения userId из токена")
		}
		return userId, nil
	}

	queryUserId := ctx.QueryInt("userId")
	if queryUserId == 0 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "Необходимо указать ?userId в QueryParams")
	}

	return uint(queryUserId), nil
}
