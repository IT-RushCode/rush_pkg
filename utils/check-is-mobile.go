package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Проверка, если запрос пришел с мобильного устройства то userID берется из JWT.
// Иначе необходимо передать userID в query.
func CheckIsMobile(ctx *fiber.Ctx) (uint, error) {
	var userId uint
	var err error

	// Безопасно извлекаем значение из контекста
	isMob, ok := ctx.Locals("IsMob").(bool)
	if !ok {
		// Если значение IsMob не найдено или не является bool
		isMob = false
	}

	if isMob {
		// Получение userId из локальных данных контекста
		userId, err = GetUserIDFromLocals(ctx)
		if err != nil {
			return 0, err
		}
	} else {
		// Получение userId из query запроса
		queryUserId := ctx.QueryInt("userId")
		if queryUserId == 0 {
			return 0, ErrorBadRequestResponse(
				ctx,
				"Необходимо указать ?userId в QueryParams",
				nil,
			)
		}
		userId = uint(queryUserId)
	}

	return userId, nil
}
