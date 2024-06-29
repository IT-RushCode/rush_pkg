package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// Проверка, если запрос пришел с мобильного устройства то userID берется из JWT.
// Иначе необходимо передать userID в query.
func CheckIsMobile(ctx *fiber.Ctx) (uint, error) {
	var userId uint
	var err error

	if ctx.QueryBool("mob") {
		// Получение userId из локальных данных контекста
		userId, err = GetUserIDFromLocals(ctx)
		if err != nil {
			return 0, err
		}
	} else {
		// Получение userId из параметров запроса
		queryUserId := ctx.QueryInt("userId")
		if queryUserId == 0 {
			return 0, errors.New("необходимо указать userId (значение не должно быть равным нулю)")

		}
		userId = uint(queryUserId)
	}

	return userId, nil
}
