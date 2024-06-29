package utils

import "github.com/gofiber/fiber/v2"

// Возвращает ID вытаскивая из params
func GetID(ctx *fiber.Ctx) (uint, error) {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return 0, ErrorBadRequestResponse(ctx, ErrorIncorrectID, nil)
	}

	return uint(id), nil
}

// Возвращает UserID вытаскивая из params
func GetUserID(ctx *fiber.Ctx) (uint, error) {
	userId, err := ctx.ParamsInt("userId")
	if err != nil {
		return 0, ErrorBadRequestResponse(ctx, ErrorIncorrectUserID, nil)
	}

	return uint(userId), nil
}
