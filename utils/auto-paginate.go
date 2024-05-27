package utils

import "github.com/gofiber/fiber/v2"

func AutoPaginate(ctx *fiber.Ctx) (limit, offset uint) {
	limit = uint(ctx.QueryInt("limit"))
	if limit == 0 {
		limit = 20
	}

	offset = uint(ctx.QueryInt("offset"))
	if offset == 0 {
		offset = 1
	}

	return limit, offset
}
