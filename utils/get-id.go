package utils

import "github.com/gofiber/fiber/v2"

// Возвращает ID вытаскивая из params
func GetID(ctx *fiber.Ctx) (uint, error) {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return 0, ErrorBadRequestResponse(ctx, "некорректный :id в параметре пути", nil)
	}

	return uint(id), nil
}
