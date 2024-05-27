package utils

import "github.com/gofiber/fiber/v2"

func ParseAndValidateInput(ctx *fiber.Ctx, input interface{}) error {
	if err := ctx.BodyParser(input); err != nil {
		return ErrorResponse(ctx, err.Error(), nil)
	}
	if err := ValidateStruct(input); err != nil {
		return ErrorResponse(ctx, err.Error(), nil)
	}
	return nil
}
